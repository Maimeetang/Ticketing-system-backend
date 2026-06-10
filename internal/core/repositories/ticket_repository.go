package repositories

import (
	"context"
	m "ticketing-system/internal/core/models"
)

type TicketRepository interface {
	Create(ctx context.Context, ticket *m.Ticket) (*m.Ticket, error)
	GetByID(ctx context.Context, id uint) (*m.Ticket, error)
	GetByCodeForUpdate(ctx context.Context, ticketcode string) (*m.Ticket, error)
	Update(ctx context.Context, ticket *m.Ticket) error
	// GetByShiftID(ctx context.Context, shiftID uint) ([]domain.Ticket, error)
}