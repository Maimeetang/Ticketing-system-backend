package port

import "ticketing-system/internal/core/domain"

type TicketTypeRepository interface {
	// CreateOrder creates a new order with its order tickets and tickets.
	CreateTicketPrice(order *domain.TicketType) (*domain.TicketType, error)
	// GetOrderByID retrieves an order with its tickets.
	GetTicketPrice(id uint) (float64, error)
	// GetOrderByUUID retrieves an order using its global UUID.
}