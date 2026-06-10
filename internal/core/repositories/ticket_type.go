package repositories

import m "ticketing-system/internal/core/models"

type TicketTypeRepository interface {
	CreateTicketType(ticketType *m.TicketType) (*m.TicketType, error)
	UpdateTicketType(ticketType *m.TicketType) (*m.TicketType, error)
	GetByID(id uint) (*m.TicketType, error)
	GetByName(name string) (*m.TicketType, error)
	ListTicketType() ([]m.TicketType, error)
}