package repositories

import (
	"context"
	m "ticketing-system/internal/core/models"
)

type TicketTypeRepository interface {
	Create(ctx context.Context, ticketType *m.TicketType) (*m.TicketType, error)
	Update(ctx context.Context, ticketType *m.TicketType) (*m.TicketType, error)
	SetActive(ctx context.Context, id uint, active bool) error
	GetByID(ctx context.Context, id uint) (*m.TicketType, error)
	List(ctx context.Context, withDisable bool) ([]m.TicketType, error)
}