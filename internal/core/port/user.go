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
	Add(user domain.User) error
	Edit(user domain.User) error
	DeleteByID(id uint) error
	FindByID(id uint) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	List() ([]domain.User, error)
}