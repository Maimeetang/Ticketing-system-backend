package port

import "ticketing-system/internal/core/domain"

type OrderService interface {
	// CreateOrder creates a new order with its order tickets and tickets.
	CreateOrder(order *domain.Order, ticketTypeID uint) (*domain.Order, error)
	// GetOrderByID retrieves an order with its tickets.
	GetOrderByID(id uint) (*domain.Order, error)
	// ListOrders returns all orders.
	ListOrders(filter domain.OrderFilter) ([]domain.Order, int64, error)
	// CancelOrder cancels an order and its tickets.
	CancelOrder(code string, userID uint) error
}

type OrderRepository interface {
	CreateOrder(order *domain.Order) (*domain.Order, error)

	GetOrderByID(id uint) (*domain.Order, error)

	GetByTicketCode(code string) (*domain.Order, error)

	ListOrders(filter domain.OrderFilter) ([]domain.Order, int64, error)

	UpdateOrder(order *domain.Order) (*domain.Order, error)
}
