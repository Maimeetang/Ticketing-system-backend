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
func (s *userServiceImpl) AddUser(user domain.User) error {
	if !validateRole(user.Role) {
		return apperror.NewBadRequest("invalid role: allowed values are cashier or scanner")
	}

	existingUser, err := s.repo.GetByUsername(user.Username)
	if err == nil && existingUser != nil {
		return apperror.NewConflict("username already exists")
	}

	hashedPassword, err := 
		bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return apperror.NewInternalServerError("failed to process security constraints: " + err.Error())
	}

	user.Password = string(hashedPassword)

	return s.repo.Create(user)
}

func (s *userServiceImpl) EditUser(user domain.User) error {
	if !validateRole(user.Role) {
		return apperror.NewBadRequest("invalid role: allowed values are cashier or scanner")
	}

	return s.repo.Update(user)
}

func (s *userServiceImpl) DeleteUser(id uint) error {
	return s.repo.DeleteByID(id)
}

func (s *userServiceImpl) FindUserByID(id uint) (*domain.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userServiceImpl) FindUserByUsername(username string) (*domain.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userServiceImpl) ListUsers() ([]domain.User, error) {
	return s.repo.GetAll()
}

// helper function
func validateRole(role domain.UserRole) bool {
	return role == domain.RoleCashier || role == domain.RoleScanner
}