package service

import (
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"
	"ticketing-system/internal/core/util"
)

type userServiceImpl struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) port.UserService{
	return &userServiceImpl{repo: repo}
}

// Business Core logic
func (s *userServiceImpl) usernameValidation(id uint, username string) error {
	existingUser, err := s.repo.GetByUsername(username)
	if err != nil {
		return err
	}

	if existingUser != nil && existingUser.ID != id {
		return apperror.NewConflict("username นี้ถูกใช้งานแล้ว")
	}

	return nil
}

func (s *userServiceImpl) RegisterUser(user *domain.User) error {
	err := s.usernameValidation(user.ID, user.Username)
	if err != nil{
		return err
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return apperror.NewInternalServerError("ไม่สามารถเข้ารหัสรหัสผ่านได้: " + err.Error())
	}

	user.Password = string(hashedPassword)

	return s.repo.Create(user)
}

func (s *userServiceImpl) UpdateUser(user *domain.User) error {
	err := s.usernameValidation(user.ID, user.Username)
	if err != nil{
		return err
	}

	return s.repo.Update(user)
}

func (s *userServiceImpl) DisableUser(id uint) error {
	return s.repo.SetActive(id, false)
}

func (s *userServiceImpl) EnableUser(id uint) error {
	return s.repo.SetActive(id, true)
}

func (s *userServiceImpl) GetUserByID(id uint) (*domain.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userServiceImpl) ListUsers() ([]domain.User, error) {
	return s.repo.List()
}