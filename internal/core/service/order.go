package service

import (
	"fmt"
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"
	"ticketing-system/internal/core/util"
)

type orderServiceImpl struct {
	orderRepo port.OrderRepository
	ticketTypeRepo port.TicketTypeRepository
}

func NewOrderService(orderRepo port.OrderRepository, ticketTypeRepo port.TicketTypeRepository) port.OrderService {
	return &orderServiceImpl{orderRepo, ticketTypeRepo}
}

func (s *orderServiceImpl) CreateOrder(order *domain.Order) (*domain.Order, error) {
	var orderTotalPrice float64

	for i := range order.Tickets {
		ticket := &order.Tickets[i]

		ticket.TicketCode = util.GenerateTicketCode(order.CashierID)
		ticket.Status = domain.TicketActive
		
		var ticketTotalPrice float64

		for j := range ticket.TicketInfos {
			info := &ticket.TicketInfos[j]

			price, err := s.ticketTypeRepo.GetTicketPrice(info.TicketTypeID)
			if err != nil {
				return nil, err
			}

			info.PricePerUnit = price

			currentSum := price * float64(info.Quantity)
			ticketTotalPrice += currentSum
			orderTotalPrice += currentSum
		}
		ticket.TotalPrice = ticketTotalPrice

		ticket.TicketLogs = append(ticket.TicketLogs, domain.TicketLog{
			FromStatus:  nil,
			ToStatus:    domain.TicketActive,
			TriggeredBy: order.CashierID,
			Remarks:     "Ticket has benn created",
		})
	}
	order.TotalPrice = orderTotalPrice
	order.Status = domain.OrderStatusPaid

	return s.orderRepo.CreateOrder(order)
}

func (s *orderServiceImpl) GetOrderByID(id uint) (*domain.Order, error) {
	return s.orderRepo.GetOrderByID(id)
}

func (s *orderServiceImpl) ListOrders() ([]domain.Order, error) {
	return s.orderRepo.ListOrders()
}

func (s *orderServiceImpl) CancelOrder(id uint, userID uint) error {
	order, err := s.orderRepo.GetOrderByID(id)
	if err != nil {
		return err
	}
	if order == nil {
		return apperror.NewNotFound("order not found.")
	}

	for _, ticket := range order.Tickets {
		if ticket.Status == domain.TicketUsed {
			return apperror.NewConflict(fmt.Sprintf("cannot cancel order: ticket %s has already been used", ticket.TicketCode))
		}
	}

	for i := range order.Tickets {
		ticket := &order.Tickets[i]
		
		oldStatus := ticket.Status 
		ticket.Status = domain.TicketCancelled

		ticket.TicketLogs = append(ticket.TicketLogs, domain.TicketLog{
			FromStatus:  &oldStatus,
			ToStatus:    ticket.Status,
			TriggeredBy: userID,
			Remarks:     "Ticket has benn cancelled",
		})
	}

	order.Status = domain.OrderStatus(domain.OrderStatusCancelled)
	
	_, err = s.orderRepo.UpdateOrder(order)
	if err != nil {
		return err
	}
	return nil
}