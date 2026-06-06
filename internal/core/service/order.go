package service

import (
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"
	"ticketing-system/internal/core/utils"
)

type orderServiceImpl struct {
	shiftRepo port.ShiftRepository
	orderRepo port.OrderRepository
	ticketTypeRepo port.TicketTypeRepository
}

func NewOrderService(shiftRepo port.ShiftRepository, orderRepo port.OrderRepository, ticketTypeRepo port.TicketTypeRepository) port.OrderService {
	return &orderServiceImpl{shiftRepo, orderRepo, ticketTypeRepo}
}

func (s *orderServiceImpl) CreateOrder(order *domain.Order, ticketTypeID uint) (*domain.Order, error) {
	shift, err := s.shiftRepo.GetCurrentByUserID(order.UserID)
	
	if err != nil {
		return nil, err
	}

	if shift == nil {
		return nil, apperror.NewForbidden("กรุณาเริ่มกะการทำงานก่อนทำรายการขาย")
	}

	order.ShiftID = shift.ID

	ticket := &order.Ticket

	ticket.TicketCode = utils.GenerateTicketCode(order.UserID)
	ticket.Status = domain.TicketActive

	var total float64

	ticketType, err := s.ticketTypeRepo.GetByID(ticketTypeID)
	if err != nil {
		return nil, err
	}

	if ticketType == nil {
		return nil, apperror.NewNotFound("ไม่พบประเภทตั๋ว")
	}

	ticket.TicketType = ticketType.Name
	ticket.PricePerUnit = ticketType.Price
	total += ticketType.Price * float64(ticket.Quantity)

	order.TotalPrice = total
	order.Status = domain.OrderStatusPaid

	ticket.TicketLogs = append(ticket.TicketLogs, domain.TicketLog{
		UserID: 	 order.UserID,
		FromStatus:  nil,
		ToStatus:    domain.TicketActive,
		Remarks:     "สร้างตั๋วสำเร็จ",
	})

	return s.orderRepo.CreateOrder(order)
}

func (s *orderServiceImpl) GetOrderByID(id uint) (*domain.Order, error) {
	return s.orderRepo.GetOrderByID(id)
}

func (s *orderServiceImpl) ListOrders(filter domain.OrderFilter) ([]domain.Order, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	return s.orderRepo.ListOrders(filter)
}

func (s *orderServiceImpl) CancelOrder(ticketCode string, userID uint) error {
	order, err := s.orderRepo.GetByTicketCode(ticketCode)
	if err != nil {
		return err
	}

	if order == nil {
		return apperror.NewNotFound("ไม่พบตั๋ว")
	}

	if order.Status == domain.OrderStatusCancelled {
		return apperror.NewConflict("รายการนี้ถูกยกเลิกแล้ว")
	}

	if order.Ticket.Status == domain.TicketUsed {
		return apperror.NewConflict("ไม่สามารถยกเลิกได้: ตั๋วถูกใช้งานแล้ว")
	}

	oldStatus := order.Ticket.Status

	order.Status = domain.OrderStatusCancelled
	order.TotalPrice = 0

	order.Ticket.Status = domain.TicketCancelled

	order.Ticket.TicketLogs = append(order.Ticket.TicketLogs, domain.TicketLog{
		UserID: 	 order.UserID,
		FromStatus:  &oldStatus,
		ToStatus:    domain.TicketCancelled,
		Remarks:     "ยกเลิกตั๋ว (พิมพ์ผิด)",
	})

	_, err = s.orderRepo.UpdateOrder(order)
	return err
}