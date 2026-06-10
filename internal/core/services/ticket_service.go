package services

import (
	"context"
	"errors"
	e "ticketing-system/internal/core/error"
	m "ticketing-system/internal/core/models"
	r "ticketing-system/internal/core/repositories"
	"time"
)

type ScanStatus string

const (
	ScanNotFound = "NOT_FOUND"
	ScanCancelled = "CANCELLED"
	ScanUsed = "USED"
	ScanSuccess = "SUCCESS"
)

type StatusUpdateResult struct {
	Status  ScanStatus
	Ticket *m.Ticket
}

type TicketService interface {
	CreateTicket(
		ctx context.Context, 
		userID uint, 
		ticketTypeID uint, 
		quantity uint,
	) (*m.Ticket, error)

	UseTicket(
		ctx context.Context, 
		userID uint, 
		ticketcode string,
	) (*StatusUpdateResult, error)

	CancelTicket(
		ctx context.Context, 
		userID uint, 
		ticketcode string, 
		remarks string,
	) (*StatusUpdateResult, error)
	
	GetTicketByID(
		ctx context.Context, id uint,
	) (*m.Ticket, error)
}

type ticketServiceImpl struct {
	transactor r.Transactor
	shiftRepo r.ShiftRepository
	ticketRepo r.TicketRepository
	ticketLogRepo r.TicketLogRepository
	ticketTypeRepo r.TicketTypeRepository
}

func NewTicketService(
	transactor r.Transactor,
	shiftRepo r.ShiftRepository,
	ticketRepo r.TicketRepository, 
	ticketLogRepo r.TicketLogRepository,
	ticketTypeRepo r.TicketTypeRepository,
) TicketService {
	return &ticketServiceImpl{
		transactor: transactor,
		shiftRepo: shiftRepo,
		ticketRepo: ticketRepo,
		ticketLogRepo: ticketLogRepo,
		ticketTypeRepo: ticketTypeRepo,
	}
}


func (s *ticketServiceImpl) CreateTicket(
	ctx context.Context,
	userID uint, 
	ticketTypeID uint, 
	quantity uint,
) (*m.Ticket, error) {
	if quantity <= 0 {
		return nil, e.NewBadRequest("quantity must be greater than 0")
	}

	var createdTicket *m.Ticket

	err := s.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		currentShift, err := s.shiftRepo.GetCurrentByUserIDForUpdate(txCtx, userID)
		if err != nil {
			return err
		}
		if currentShift == nil {
			return e.NewForbidden("please open a shift before selling tickets")
		}

		ticketType, err := s.ticketTypeRepo.GetByID(ticketTypeID)
		if err != nil {
			return err
		}
		if ticketType == nil {
			return e.NewBadRequest("invalid ticket type id")
		}

		ticketCode := generateTicketCode(userID)
		unitPrice := ticketType.Price
		totalPrice := unitPrice * int64(quantity)

		now := time.Now()
		
		ticket := &m.Ticket{
			ShiftID: currentShift.ID,
			TicketTypeID: ticketTypeID,
			TicketTypeName: ticketType.Name,
			TicketCode: ticketCode,
			Quantity: quantity,
			UnitPrice: unitPrice,
			TotalPrice: totalPrice,
			Status: m.TicketActive,
			SoldAt: now,
		}

		newTicket, err := s.ticketRepo.Create(txCtx, ticket)
		if err != nil {
			return err
		}

		createdTicket = newTicket
		return nil
	})
	if err != nil {
		return nil, err
	}

	return createdTicket, nil
}

func (s *ticketServiceImpl) GetTicketByID(
	ctx context.Context, 
	id uint,
) (*m.Ticket, error) {
	return s.ticketRepo.GetByID(ctx, id)
}

// Private Helper Transaction Function Service
func (s *ticketServiceImpl) executeStatusChangeTx(
	ctx context.Context, 
	ticketCode string, 
	userID uint,
	targetStatus m.TicketStatus,
	remarks string,
	validateFunc func(ctx context.Context, ticket *m.Ticket) error,
) (*m.Ticket, error) {
	var updatedTicket *m.Ticket

	err := s.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		ticket, err := s.ticketRepo.GetByCodeForUpdate(txCtx, ticketCode)
		if err != nil {
			return err
		}
		if ticket == nil {
			return errors.New(ScanNotFound)
		}

		if err := validateFunc(txCtx, ticket); err != nil {
			return err
		}

		// update data
		fromStatus := ticket.Status
		ticket.Status = targetStatus

		now := time.Now()

		switch targetStatus {
		case m.TicketUsed:
			ticket.UsedAt = &now
		case m.TicketCancelled:
			ticket.CancelledAt = &now
		}


		log := &m.TicketLog{
			UserID:     userID,
			TicketID:   ticket.ID,
			FromStatus: fromStatus,
			ToStatus:   targetStatus,
			Remarks:    remarks,
		}

		// save data
		if err := s.ticketRepo.Update(txCtx, ticket); err != nil {
			return err
		}
		if err := s.ticketLogRepo.Create(txCtx, log); err != nil {
			return err
		}

		updatedTicket = ticket
		return nil
	})

	if err != nil {
		return nil, err
	}
	return updatedTicket, nil
}

func (s *ticketServiceImpl) UseTicket(
	ctx context.Context, 
	userID uint, 
	code string,
) (*StatusUpdateResult, error) {
	ticket, err := s.executeStatusChangeTx(
		ctx, code, userID, m.TicketUsed, "", func(txCtx context.Context, t *m.Ticket,
		) error {
		if t.Status == m.TicketCancelled {
			return errors.New(ScanCancelled)
		}
		if t.Status == m.TicketUsed {
			return errors.New(ScanUsed)
		}
		return nil
	})

	if err != nil {
		switch err.Error() {
		case ScanNotFound:
			return &StatusUpdateResult{Status: ScanNotFound}, nil
		case ScanCancelled:
			return &StatusUpdateResult{Status: ScanCancelled}, nil
		case ScanUsed:
			return &StatusUpdateResult{Status: ScanUsed}, nil
		default:
			return nil, err
		}
	}

	return &StatusUpdateResult{Status: ScanSuccess, Ticket: ticket}, nil
}

func (s *ticketServiceImpl) CancelTicket(
	ctx context.Context, 
	userID uint, 
	code string, 
	remarks string,
) (*StatusUpdateResult, error) {
	ticket, err := s.executeStatusChangeTx(
		ctx, code, userID, m.TicketCancelled, remarks, func(txCtx context.Context, t *m.Ticket,
		) error {
		if t.Status == m.TicketCancelled {
			return errors.New(ScanCancelled)
		}
		if t.Status == m.TicketUsed {
			return errors.New(ScanUsed)
		}
		shift, err := s.shiftRepo.GetByIDForUpdate(txCtx, t.ShiftID)
		if err != nil {
			return err
		}
		if shift != nil && shift.Status == m.ShiftClosed {
			return e.NewConflict("cannot cancel ticket because the shift is already closed")
		}
		return nil
	})

	if err != nil {
		switch err.Error() {
		case ScanNotFound:
			return &StatusUpdateResult{Status: ScanNotFound}, nil
		case ScanCancelled:
			return &StatusUpdateResult{Status: ScanCancelled}, nil
		case ScanUsed:
			return &StatusUpdateResult{Status: ScanUsed}, nil
		default:
			return nil, err
		}
	}

	return &StatusUpdateResult{Status: ScanSuccess, Ticket: ticket}, nil
}