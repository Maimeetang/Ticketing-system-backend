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
	Create(user *domain.User) error
	Update(user *domain.User) error
	SetActive(id uint, active bool) error
	GetByID(id uint) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	List() ([]domain.User, error)
}