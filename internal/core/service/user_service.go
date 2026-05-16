package service

import (
	coreerr "ticketing-system/internal/core/errors"
	"ticketing-system/internal/core/model"
	"ticketing-system/internal/core/repository"
)

type UserService interface {
	AddUser(user model.User) error
	EditUser(user model.User) error
	DeleteUser(ID uint) error
	FindUserByID(ID uint) (model.User, error)
	FindUserByUsername(username string) (model.User, error)
}

type userServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService{
	return &userServiceImpl{repo: repo}
}

func validateRole(role string) bool {
	return role == "manager" || role == "cashier"
}

// Business Core logic
func (s *userServiceImpl) AddUser(user model.User) error {
	if !validateRole(user.Role) {
		return coreerr.NewBadRequest("invalid role")
	}
	return s.repo.Add(user)
}

func (s *userServiceImpl) EditUser(user model.User) error {
	return s.repo.Edit(user)
}

func (s *userServiceImpl) DeleteUser(id uint,) error {
	return s.repo.DeleteByID(id)
}

func (s *userServiceImpl) FindUserByID(id uint) (model.User, error) {
	return s.repo.FindByID(id)
}

func (s *userServiceImpl) FindUserByUsername(username string) (model.User, error) {
	return s.repo.FindByUsername(username)
}