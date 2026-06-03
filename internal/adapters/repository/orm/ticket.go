package orm

import (
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"

	"gorm.io/gorm"
)

type GormTicketRepository struct {
	db *gorm.DB
}

func NewGormTicketRepository(db *gorm.DB) port.TicketRepository {
	return &GormTicketRepository{db: db}
}

func (r *GormTicketRepository) GetByCode(code string) (*domain.Ticket, error) {
	var ticket domain.Ticket
	err := r.db.
		Preload("TicketLog").
		Preload("TicketType").
		Where("ticket_code = ?", code).
		First(&ticket).Error

	if err != nil {
		return nil, handleError(err)
	}

	return &ticket, nil
}

func (r *GormTicketRepository) UpdateTicket(ticket *domain.Ticket) error {
	err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&ticket).Error
	if err != nil {
		return handleError(err)
	}

	return nil
}