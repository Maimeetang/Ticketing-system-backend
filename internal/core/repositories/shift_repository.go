package repositories

import (
	"context"
	m "ticketing-system/internal/core/models"
	"time"
)

type ShiftRepository interface {
	Create(ctx context.Context, shift *m.Shift) error
	Update(ctx context.Context, shift *m.Shift) error
	GetByID(ctx context.Context, id uint) (*m.Shift, error)
	GetByIDForUpdate(ctx context.Context, id uint) (*m.Shift, error)
	GetCurrentByUserID(ctx context.Context, userID uint) (*m.Shift, error)
	GetCurrentByUserIDForUpdate(ctx context.Context, userID uint) (*m.Shift, error)
	GetByDate(ctx context.Context, date time.Time) ([]m.Shift, error)
	CalculateSummary(ctx context.Context, shift *m.Shift) error
}
