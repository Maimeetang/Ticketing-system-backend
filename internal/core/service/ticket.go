package service

import (
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"
)

type TicketServiceImpl struct {
	repo port.TicketRepository
}

func NewTicketService(repo port.TicketRepository) port.TicketService {
	return &TicketServiceImpl{repo: repo}
}

func (s *TicketServiceImpl) UseTicket(code string, userID uint) error {
	ticket, err := s.repo.GetByCode(code)

	if err != nil {
		return err
	}

	if ticket == nil {
		return apperror.NewNotFound("ไม่พบข้อมูลตั๋ว")
	}

	if ticket.Status == domain.TicketUsed {
		return apperror.NewConflict("ตั๋วนี้ถูกใช้งานแล้ว")
	}

	if ticket.Status == domain.TicketCancelled {
		return apperror.NewConflict("ตั๋วนี้ถูกยกเลิกแล้ว")
	}

	oldStatus := ticket.Status
	ticket.Status = domain.TicketUsed

	ticket.TicketLogs = append(ticket.TicketLogs, domain.TicketLog{
		FromStatus:  &oldStatus,
		ToStatus:    domain.TicketUsed,
		TriggeredBy: userID,
		Remarks:     "สแกนใช้งานตั๋ว",
	})

	return s.repo.UpdateTicket(ticket)
}