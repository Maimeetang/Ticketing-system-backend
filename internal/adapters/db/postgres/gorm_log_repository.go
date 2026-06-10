package orm

import (
	"context"
	m "ticketing-system/internal/core/models"
	r "ticketing-system/internal/core/repositories"

	"gorm.io/gorm"
)
type GormTicketLogRepository struct {
	db *gorm.DB
}

func NewgormTicketLogRepository(db *gorm.DB) r.TicketLogRepository {
	return &GormTicketLogRepository{db: db}
}

func (r *GormTicketLogRepository) Create(ctx context.Context, log *m.TicketLog) error {
	db := bind(ctx, r.db)

	return db.Create(log).Error
}