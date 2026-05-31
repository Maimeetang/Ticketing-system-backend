package service

import (
	"strings"
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"
)

type ticketTypeServiceImpl struct {
	repo port.TicketTypeRepository
}

func NewTicketTypeService(repo port.TicketTypeRepository) port.TicketTypeService {
	return &ticketTypeServiceImpl{repo: repo}
}

// Business Core logic
func (s *ticketTypeServiceImpl) validateAndNormalizeName(ticketType *domain.TicketType) error {
	ticketType.Name = strings.ToUpper(ticketType.Name)

	existing, err := s.repo.GetByName(ticketType.Name)
	if err != nil {
		return err
	}

	if existing != nil && existing.ID != ticketType.ID {
		return apperror.NewConflict("มีประเภทตั๋วนี้อยู่แล้ว")
	}

	return nil
}

func (s *ticketTypeServiceImpl) CreateTicketType(ticketType *domain.TicketType) (*domain.TicketType, error) {
	if err := s.validateAndNormalizeName(ticketType); err != nil {
		return nil, err
	}

	return s.repo.CreateTicketType(ticketType)
}

func (s *ticketTypeServiceImpl) UpdateTicketType(ticketType *domain.TicketType) (*domain.TicketType, error) {
	if err := s.validateAndNormalizeName(ticketType); err != nil {
		return nil, err
	}

	return s.repo.UpdateTicketType(ticketType)
}

func (s *ticketTypeServiceImpl) GetTicketType(id uint) (*domain.TicketType, error) {
	return s.repo.GetByID(id)
}

func (s *ticketTypeServiceImpl) ListTicketType() ([]domain.TicketType, error) {
	return s.repo.ListTicketType()
}