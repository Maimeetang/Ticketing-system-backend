package repositories

import (
	"context"
	m "ticketing-system/internal/core/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *m.User) error
	Update(ctx context.Context, user *m.User) error
	SetActive(ctx context.Context, id uint, active bool) error
	GetByID(ctx context.Context, id uint) (*m.User, error)
	GetByIDForUpdate(ctx context.Context, id uint) (*m.User, error)
	GetByUsername(ctx context.Context, username string) (*m.User, error)
	List(ctx context.Context, ) ([]m.User, error)
}