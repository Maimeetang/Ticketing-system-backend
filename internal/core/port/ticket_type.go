package port

import "ticketing-system/internal/core/domain"

type TicketTypeService interface {
	CreateTicketType(ticketType *domain.TicketType) (*domain.TicketType, error)
	UpdateTicketType(ticketType *domain.TicketType) (*domain.TicketType, error)
	GetTicketType(id uint) (*domain.TicketType, error)
	ListTicketType() ([]domain.TicketType, error)
}

type TicketTypeRepository interface {
	CreateTicketType(ticketType *domain.TicketType) (*domain.TicketType, error)
	UpdateTicketType(ticketType *domain.TicketType) (*domain.TicketType, error)
	GetByID(id uint) (*domain.TicketType, error)
	GetByName(name string) (*domain.TicketType, error)
	ListTicketType() ([]domain.TicketType, error)
}