package services

import (
	"context"
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
	EnableTicketType(ctx context.Context, id uint) error
	DisableTicketType(ctx context.Context, id uint) error
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
	Ttype := &m.TicketType{
		ID: id,
		Name: name,
		Price: price,
		Description: desc,
	}
	return s.repo.Update(ctx, Ttype)
}

func (s *ticketTypeServiceImpl) GetTicketType(ctx context.Context, id uint) (*m.TicketType, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ticketTypeServiceImpl) ListTicketType(
	ctx context.Context, withDisable bool,
) ([]m.TicketType, error) {
	return s.repo.List(ctx, withDisable)
}

func (s *ticketTypeServiceImpl) EnableTicketType(ctx context.Context, id uint) error {
	return s.repo.SetActive(ctx, id, true)
}

func (s *ticketTypeServiceImpl) DisableTicketType(ctx context.Context, id uint) error {
	return s.repo.SetActive(ctx, id, false)
}