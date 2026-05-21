package port

import "ticketing-system/internal/core/domain"

// Primary port
type UserService interface {
	// CreateUser registers a new employee
	Register(user *domain.User) error
	// UpdateUser updates a user
	UpdateUser(user *domain.User) error
	// DeleteUser deletes a user
	DeleteUser(ID uint) error
	// GetUser returns a user by id
	GetUser(ID uint) (*domain.User, error)
	// ListUsers return a list of all users
	ListUsers() ([]domain.User, error)
}

// Secondary port
type UserRepository interface {
	// CreateUser inserts a new user into the database
	CreateUser(user *domain.User) error
	// UpdateUser updates a user
	UpdateUser(user *domain.User) error
	// DeleteUser deletes a user
	DeleteUser(id uint) error
	// GetUserByID selects a user by id
	GetUserByID(id uint) (*domain.User, error)
	// GetByUsername selects a user by username
	GetUserByUsername(username string) (*domain.User, error)
	// ListUsers selects a list of all users
	ListUsers() ([]domain.User, error)
}