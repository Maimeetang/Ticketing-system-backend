package orm

import (
	"context"
	domains "ticketing-system/internal/core/models"
	r "ticketing-system/internal/core/repositories"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormTicketRepository struct {
	db *gorm.DB
}

func NewGormTicketRepository(db *gorm.DB) r.TicketRepository {
	return &gormTicketRepository{db: db}
}

func (r *gormTicketRepository) Create(
	ctx context.Context,
	ticket *domains.Ticket,
) (*domains.Ticket, error) {
	db := bind(ctx, r.db)

	if err := db.Create(ticket).Error; err != nil {
		return nil, handleError(err)
	}
	return ticket, nil
}

func (r *gormTicketRepository) GetByID(
	ctx context.Context,
	id uint,
) (*domains.Ticket, error) {
	var ticket domains.Ticket
	db := bind(ctx, r.db)

	err := db.WithContext(ctx).First(&ticket, id).Error

	return &ticket, handleError(err)
}

func (r *gormTicketRepository) GetByCodeForUpdate(
	ctx context.Context,
	code string,
) (*domains.Ticket, error) {
	var ticket domains.Ticket
	db := bind(ctx, r.db)

	err := db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("ticket_code = ?", code).
		First(&ticket).Error

	return &ticket, handleError(err)
}

func (r *gormTicketRepository) Update(
	ctx context.Context, 
	ticket *domains.Ticket,
) error {
	db := bind(ctx, r.db)

	return db.Updates(ticket).Error
}