package port

import "ticketing-system/internal/core/domain"

// Primary port
type UserService interface {
	AddUser(user domain.User) error
	EditUser(user domain.User) error
	DeleteUser(ID uint) error
	FindUserByID(ID uint) (*domain.User, error)
	FindUserByUsername(username string) (*domain.User, error)
	ListUsers() ([]domain.User, error)
}

// Secondary port
type UserRepository interface {
	Create(user domain.User) error
	Update(user domain.User) error
	DeleteByID(id uint) error
	GetByID(id uint) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	GetAll() ([]domain.User, error)
}