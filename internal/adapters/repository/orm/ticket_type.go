package orm

import (
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"

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
		return nil, handleTicketTypeError(err)
	}
	return ticketType, nil
}

func (r *GormTicketTypeRepository) GetByID(id uint) (*domain.TicketType, error) {
	var ticketType domain.TicketType
	err := r.db.First(&ticketType, id).Error
	if err != nil {
		return nil, handleTicketTypeError(err)
	}
	return &ticketType, nil
}

func (r *GormTicketTypeRepository) GetByName(name string) (*domain.TicketType, error) {
	var ticketType domain.TicketType
	err := r.db.Where("name = ?", name).First(&ticketType).Error
	if err != nil {
		return nil, handleTicketTypeError(err)
	}
	return &ticketType, nil
}

func (r *GormTicketTypeRepository) UpdateTicketType(ticketType *domain.TicketType) (*domain.TicketType, error) {
	result := r.db.Model(&domain.TicketType{}).
		Where("id = ?", ticketType.ID).
		Select("*").
		Updates(ticketType)

	if result.Error != nil {
		return nil, handleTicketTypeError(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, apperror.NewNotFound("ticket type")
	}

	var updated domain.TicketType
	err := r.db.First(&updated, ticketType.ID).Error
	if err != nil {
		return nil, handleTicketTypeError(err)
	}

	return &updated, nil
}

func handleTicketTypeError(err error) error {
	if err == nil {
		return nil
	}
	return apperror.NewInternalServerError("database error: " + err.Error())
}