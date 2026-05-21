package port

import "ticketing-system/internal/core/domain"

type OrderService interface {
	CreateOrder(userID uint, shiftID uint, paymentMethod domain.PaymentMethod, items []domain.OrderItem) (*domain.Order, error)
}

type OrderRepository interface {
	Create(order *domain.Order) error
	GetTicketTypeByID(id uint) (*domain.TicketType, error)
}
