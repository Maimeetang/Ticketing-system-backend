package port

import "ticketing-system/internal/core/domain"

type UserService interface {
	RegisterUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	DisableUser(ID uint) error
	EnableUser(ID uint) error
	GetUserByID(ID uint) (*domain.User, error)
	ListUsers() ([]domain.User, error)
}

type UserRepository interface {
	CreateUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	SetActiveUser(id uint, active bool) error
	GetUserByID(id uint) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	ListUsers() ([]domain.User, error)
}