package db

import (
	"context"
	e "ticketing-system/internal/core/error"
	m "ticketing-system/internal/core/models"
	r "ticketing-system/internal/core/repositories"

	"gorm.io/gorm"
)

type GormTicketTypeRepository struct {
	db *gorm.DB
}

func NewGormTicketTypeRepository(db *gorm.DB) r.TicketTypeRepository {
	return &GormTicketTypeRepository{db: db}
}

func (r *GormTicketTypeRepository) Create(
	ctx context.Context, ticketType *m.TicketType,
) (*m.TicketType, error) {
	db := bind(ctx, r.db)

	if err := db.Create(ticketType).Error; err != nil{ 
		return nil, handleError(err)
	}
	return ticketType, nil
}

func (r *GormTicketTypeRepository) GetByID(
	ctx context.Context, id uint,
) (*m.TicketType, error) {
	var ticketType m.TicketType
	db := bind(ctx, r.db)

	if err := db.First(&ticketType, id).Error; err!= nil {
		return nil, handleError(err)
	}
	return &ticketType, nil
}

func (r *GormTicketTypeRepository) SetActive(ctx context.Context, id uint, active bool) error {
	db := bind(ctx, r.db)

	result := db.
		Model(&m.User{}).
		Where("id = ?", id).
		Update("is_active", active)

	if result.RowsAffected == 0 {
		return e.NewNotFound("ticket type not found")
	}

	return handleError(result.Error)
}

func (r *GormTicketTypeRepository) Update(
	ctx context.Context, ticketType *m.TicketType,
) (*m.TicketType, error) {
	db := bind(ctx, r.db)

	if err := db.Updates(ticketType).Error; err != nil {
		return nil, handleError(err)
	}
	return ticketType, nil
}

func (r *GormTicketTypeRepository) List(ctx context.Context, withDisable bool) ([]m.TicketType, error) {
	var types []m.TicketType
	db := bind(ctx, r.db)

	if !withDisable {
		db = db.Where("is_active = ?", true)
	}

	if err := db.Find(&types).Error; err != nil {
		return nil, handleError(err)
	}
	return types, nil
}