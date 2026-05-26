package service

import (
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"

	"golang.org/x/crypto/bcrypt"
)

type userServiceImpl struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) port.UserService{
	return &userServiceImpl{repo: repo}
}

// Business Core logic
func (s *userServiceImpl) Register(user *domain.User) error {
	existingUser, err := s.repo.GetUserByUsername(user.Username)
	if err == nil && existingUser != nil {
		return apperror.NewConflict("username already exists")
	}

	hashedPassword, err := 
		bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return apperror.NewInternalServerError("failed to process security constraints: " + err.Error())
	}

	user.Password = string(hashedPassword)

	return s.repo.CreateUser(user)
}

func (s *userServiceImpl) UpdateUser(user *domain.User) error {
	return s.repo.UpdateUser(user)
}

func (s *userServiceImpl) DeleteUser(id uint) error {
	return s.repo.DeleteUser(id)
}

func (s *userServiceImpl) GetUser(id uint) (*domain.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userServiceImpl) ListUsers() ([]domain.User, error) {
	return s.repo.ListUsers()
}