package services

import (
	"context"
	e "ticketing-system/internal/core/error"
	m "ticketing-system/internal/core/models"
	r "ticketing-system/internal/core/repositories"
)

type TicketTypeService interface {
	CreateTicketType(
		ctx context.Context, 
		name string, 
		price int64, 
		desc string,
	) (*m.TicketType, error)
	UpdateTicketType(
		ctx context.Context, 
		id uint, 
		name string, 
		price int64, 
		desc string,
	) (*m.TicketType, error)
	UpdateTicketTypeStatus(ctx context.Context, id uint, isActive bool) error
	GetTicketType(ctx context.Context, id uint) (*m.TicketType, error)
	ListTicketType(ctx context.Context, withDisable bool) ([]m.TicketType, error)
}

type ticketTypeServiceImpl struct {
	repo r.TicketTypeRepository
}

func NewTicketTypeService(repo r.TicketTypeRepository) TicketTypeService {
	return &ticketTypeServiceImpl{repo: repo}
}

func (s *ticketTypeServiceImpl) CreateTicketType(
	ctx context.Context,
	name string, 
	price int64, 
	desc string,
) (*m.TicketType, error) {
	if price < 0 {
		return nil, e.NewBadRequest("price must be greater than 0")
	}

	Ttype := &m.TicketType{
		Name: name,
		Price: price,
		Description: desc,
		IsActive: true,
	}

	return s.repo.Create(ctx, Ttype)
}

func (s *ticketTypeServiceImpl) UpdateTicketType(
	ctx context.Context,
	id uint, 
	name string, 
	price int64, 
	desc string,
) (*m.TicketType, error) {
	if price < 0 {
		return nil, e.NewBadRequest("price must be greater than 0")
	}

	Ttype := &m.TicketType{
		ID: id,
		Name: name,
		Price: price,
		Description: desc,
	}
	return s.repo.Update(ctx, Ttype)
}

func (s *ticketTypeServiceImpl) GetTicketType(ctx context.Context, id uint) (*m.TicketType, error) {
	tType, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if tType == nil {
		return nil, e.NewNotFound("ticket type not found")
	}

	return tType, nil
}

func (s *ticketTypeServiceImpl) ListTicketType(
	ctx context.Context, withDisable bool,
) ([]m.TicketType, error) {
	return s.repo.List(ctx, withDisable)
}

func (s *ticketTypeServiceImpl) UpdateTicketTypeStatus(ctx context.Context, id uint, isActive bool) error {
	return s.repo.SetActive(ctx, id, isActive)
}