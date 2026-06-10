package services

import (
	"strings"
	e "ticketing-system/internal/core/error"
	m "ticketing-system/internal/core/models"
	r "ticketing-system/internal/core/repositories"
)

type TicketTypeService interface {
	CreateTicketType(ticketType *m.TicketType) (*m.TicketType, error)
	UpdateTicketType(ticketType *m.TicketType) (*m.TicketType, error)
	GetTicketType(id uint) (*m.TicketType, error)
	ListTicketType() ([]m.TicketType, error)
}

type ticketTypeServiceImpl struct {
	repo r.TicketTypeRepository
}

func NewTicketTypeService(repo r.TicketTypeRepository) TicketTypeService {
	return &ticketTypeServiceImpl{repo: repo}
}

func (s *ticketTypeServiceImpl) validateAndNormalizeName(ticketType *m.TicketType) error {
	ticketType.Name = strings.ToUpper(ticketType.Name)

	existing, err := s.repo.GetByName(ticketType.Name)
	if err != nil {
		return err
	}

	if existing != nil && existing.ID != ticketType.ID {
		return e.NewConflict("มีประเภทตั๋วนี้อยู่แล้ว")
	}

	return nil
}

func (s *ticketTypeServiceImpl) CreateTicketType(ticketType *m.TicketType) (*m.TicketType, error) {
	if err := s.validateAndNormalizeName(ticketType); err != nil {
		return nil, err
	}

	return s.repo.CreateTicketType(ticketType)
}

func (s *ticketTypeServiceImpl) UpdateTicketType(ticketType *m.TicketType) (*m.TicketType, error) {
	if err := s.validateAndNormalizeName(ticketType); err != nil {
		return nil, err
	}

	return s.repo.UpdateTicketType(ticketType)
}

func (s *ticketTypeServiceImpl) GetTicketType(id uint) (*m.TicketType, error) {
	return s.repo.GetByID(id)
}

func (s *ticketTypeServiceImpl) ListTicketType() ([]m.TicketType, error) {
	return s.repo.ListTicketType()
}