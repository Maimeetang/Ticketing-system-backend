package port

import "ticketing-system/internal/core/domain"

type TicketRepository interface {
	UpdateTicket(ticket *domain.Ticket) (*domain.Ticket, error)
}