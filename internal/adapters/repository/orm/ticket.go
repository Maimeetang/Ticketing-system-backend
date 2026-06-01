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
		Preload("TicketLogs").
		Where("ticket_code = ?", code).
		First(&ticket).Error

	if err != nil {
		return nil, handleError(err)
	}

	return &ticket, nil
}

func (r *GormTicketRepository) UpdateTicket(ticket *domain.Ticket) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.Ticket{}).
			Where("id = ?", ticket.ID).
			Updates(map[string]any{
				"status":      ticket.Status,
				"total_price": ticket.TotalPrice,
			}).Error; err != nil {
			return err
		}

		var newLogs []domain.TicketLog
		for i := range ticket.TicketLogs {
			log := ticket.TicketLogs[i]

			if log.ID == 0 {
				log.TicketID = ticket.ID
				newLogs = append(newLogs, log)
			}
		}

		if len(newLogs) > 0 {
			if err := tx.Create(&newLogs).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return handleError(err)
	}

	return nil
}