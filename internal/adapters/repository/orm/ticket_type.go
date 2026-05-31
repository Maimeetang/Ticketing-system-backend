package orm

import (
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"
	"time"

	"gorm.io/gorm"
)

type GormTicketTypeRepository struct {
	db *gorm.DB
}

func NewGormTicketTypeRepository(db *gorm.DB) port.TicketTypeRepository {
	return &GormTicketTypeRepository{db: db}
}

func (r *GormTicketTypeRepository) CreateTicketType(ticketType *domain.TicketType) (*domain.TicketType, error) {
	err := r.db.Create(ticketType).Error
	if err != nil {
		return nil, handleError(err)
	}
	return ticketType, nil
}

func (r *GormTicketTypeRepository) GetByID(id uint) (*domain.TicketType, error) {
	var ticketType domain.TicketType
	err := r.db.First(&ticketType, id).Error
	if err != nil {
		return nil, handleError(err)
	}
	return &ticketType, nil
}

func (r *GormTicketTypeRepository) GetByName(name string) (*domain.TicketType, error) {
	var ticketType domain.TicketType
	err := r.db.Where("name = ?", name).First(&ticketType).Error
	if err != nil {
		return nil, handleError(err)
	}
	return &ticketType, nil
}

func (r *GormTicketTypeRepository) UpdateTicketType(ticketType *domain.TicketType) (*domain.TicketType, error) {
	update := map[string]any{
		"name":        ticketType.Name,
		"price":       ticketType.Price,
		"description": ticketType.Description,
		"is_active":   ticketType.IsActive,
		"updated_at":  time.Now(),
	}

	result := r.db.Model(&domain.TicketType{}).
		Where("id = ?", ticketType.ID).
		Updates(update)

	if result.Error != nil {
		return nil, handleError(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, apperror.NewNotFound("ticket type")
	}

	var updated domain.TicketType
	err := r.db.First(&updated, ticketType.ID).Error
	if err != nil {
		return nil, handleError(err)
	}

	return &updated, nil
}

func (r *GormTicketTypeRepository) ListTicketType() ([]domain.TicketType, error) {
	var types []domain.TicketType
	err := r.db.Find(&types).Error
	if err != nil {
		return nil, handleError(err)
	}
	return types, nil
}