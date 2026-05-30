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
		return apperror.NewConflict("ticket type name already exists")
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

func (s *ticketTypeServiceImpl) DisableTicketType(id uint) error {
	ticketType, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if !ticketType.IsActive {
		return nil
	}

	ticketType.IsActive = false
	_, err = s.repo.UpdateTicketType(ticketType)
	return err
}

func (s *ticketTypeServiceImpl) GetTicketType(id uint) (*domain.TicketType, error) {
	return s.repo.GetByID(id)
}