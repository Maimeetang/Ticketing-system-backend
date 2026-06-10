package repositories

import (
	"context"
	m "ticketing-system/internal/core/models"
)

type TicketLogRepository interface {
	Create(ctx context.Context, log *m.TicketLog) error
}