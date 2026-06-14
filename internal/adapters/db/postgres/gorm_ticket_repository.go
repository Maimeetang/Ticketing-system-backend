package db

import (
	"context"
	m "ticketing-system/internal/core/models"
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
	ticket *m.Ticket,
) (*m.Ticket, error) {
	db := bind(ctx, r.db)

	if err := db.Create(ticket).Error; err != nil {
		return nil, handleError(err)
	}
	return ticket, nil
}

func (r *gormTicketRepository) GetByID(
	ctx context.Context,
	id uint,
) (*m.Ticket, error) {
	var ticket m.Ticket
	db := bind(ctx, r.db)

	err := db.First(&ticket, id).Error

	if err != nil {
        return nil, handleError(err)
    }

	return &ticket, nil
}

func (r *gormTicketRepository) GetByShiftID(
	ctx context.Context,
	shiftid uint,
) ([]m.Ticket, error) {
	var tickets []m.Ticket
	db := bind(ctx, r.db)

	err := db.Find(&tickets, "shift_id = ?", shiftid).Error

	if err != nil {
        return nil, handleError(err)
    }

	return tickets, nil
}

func (r *gormTicketRepository) GetByCodeForUpdate(
	ctx context.Context,
	code string,
) (*m.Ticket, error) {
	var ticket m.Ticket
	db := bind(ctx, r.db)

	err := db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("ticket_code = ?", code).
		First(&ticket).Error

	if err != nil {
        return nil, handleError(err)
    }

	return &ticket, nil
}

func (r *gormTicketRepository) Update(
	ctx context.Context, 
	ticket *m.Ticket,
) error {
	db := bind(ctx, r.db)

	return db.Updates(ticket).Error
}