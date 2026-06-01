package port

import "ticketing-system/internal/core/domain"

type TicketService interface {
	UseTicket(code string, userID uint) error
}

type TicketRepository interface {
	UpdateTicket(ticket *domain.Ticket) error
	GetByCode(code string) (*domain.Ticket, error)
}