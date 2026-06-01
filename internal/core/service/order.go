package service

import (
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"
	"ticketing-system/internal/core/util"
)

type orderServiceImpl struct {
	shiftRepo port.ShiftRepository
	orderRepo port.OrderRepository
	ticketTypeRepo port.TicketTypeRepository
}

func NewOrderService(shiftRepo port.ShiftRepository, orderRepo port.OrderRepository, ticketTypeRepo port.TicketTypeRepository) port.OrderService {
	return &orderServiceImpl{shiftRepo, orderRepo, ticketTypeRepo}
}

func (s *orderServiceImpl) CreateOrder(order *domain.Order) (*domain.Order, error) {
	shift, err := s.shiftRepo.GetActiveByUserID(order.CashierID)
	
	if err != nil {
		return nil, err
	}

	if shift == nil {
		return nil, apperror.NewForbidden("กรุณาเริ่มกะการทำงานก่อนทำรายการขาย")
	}

	order.ShiftID = shift.ID

	ticket := &order.Ticket

	ticket.TicketCode = util.GenerateTicketCode(order.CashierID)
	ticket.Status = domain.TicketActive

	var total float64

	for i := range ticket.TicketInfos {
		info := &ticket.TicketInfos[i]

		ticketType, err := s.ticketTypeRepo.GetByID(info.TicketTypeID)
		if err != nil {
			return nil, err
		}

		if ticketType == nil {
			return nil, apperror.NewNotFound("ไม่พบประเภทตั๋ว")
		}

		info.PricePerUnit = ticketType.Price

		currentSum := ticketType.Price * float64(info.Quantity)
		total += currentSum
	}

	ticket.TotalPrice = total
	order.TotalPrice = total
	order.Status = domain.OrderStatusPaid

	ticket.TicketLogs = append(ticket.TicketLogs, domain.TicketLog{
		FromStatus:  nil,
		ToStatus:    domain.TicketActive,
		TriggeredBy: order.CashierID,
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
		FromStatus:  &oldStatus,
		ToStatus:    domain.TicketCancelled,
		TriggeredBy: userID,
		Remarks:     "ยกเลิกตั๋ว (พิมพ์ผิด)",
	})

	_, err = s.orderRepo.UpdateOrder(order)
	return err
}