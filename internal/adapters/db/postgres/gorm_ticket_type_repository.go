package orm

import (
	e "ticketing-system/internal/core/error"
	m "ticketing-system/internal/core/models"
	r "ticketing-system/internal/core/repositories"
	"time"

	"gorm.io/gorm"
)

type GormTicketTypeRepository struct {
	db *gorm.DB
}

func NewGormTicketTypeRepository(db *gorm.DB) r.TicketTypeRepository {
	return &GormTicketTypeRepository{db: db}
}

func (r *GormTicketTypeRepository) CreateTicketType(ticketType *m.TicketType) (*m.TicketType, error) {
	err := r.db.Create(ticketType).Error
	if err != nil {
		return nil, handleError(err)
	}
	return ticketType, nil
}

func (r *GormTicketTypeRepository) GetByID(id uint) (*m.TicketType, error) {
	var ticketType m.TicketType
	err := r.db.First(&ticketType, id).Error
	if err != nil {
		return nil, handleError(err)
	}
	return &ticketType, nil
}

func (r *GormTicketTypeRepository) GetByName(name string) (*m.TicketType, error) {
	var ticketType m.TicketType
	err := r.db.Where("name = ?", name).First(&ticketType).Error
	if err != nil {
		return nil, handleError(err)
	}
	return &ticketType, nil
}

func (r *GormTicketTypeRepository) UpdateTicketType(ticketType *m.TicketType) (*m.TicketType, error) {
	update := map[string]any{
		"name":        ticketType.Name,
		"price":       ticketType.Price,
		"description": ticketType.Description,
		"is_active":   ticketType.IsActive,
		"updated_at":  time.Now(),
	}

	result := r.db.Model(&m.TicketType{}).
		Where("id = ?", ticketType.ID).
		Updates(update)

	if result.Error != nil {
		return nil, handleError(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, e.NewNotFound("ticket type")
	}

	var updated m.TicketType
	err := r.db.First(&updated, ticketType.ID).Error
	if err != nil {
		return nil, handleError(err)
	}

	return &updated, nil
}

func (r *GormTicketTypeRepository) ListTicketType() ([]m.TicketType, error) {
	var types []m.TicketType
	err := r.db.Find(&types).Error
	if err != nil {
		return nil, handleError(err)
	}
	return types, nil
}