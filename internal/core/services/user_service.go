package services

import (
	"context"
	e "ticketing-system/internal/core/error"
	m "ticketing-system/internal/core/models"
	r "ticketing-system/internal/core/repositories"
)

type UserService interface {
	RegisterUser(
		ctx context.Context, 
		username string,
		password string,
		role string,
		firstName string,
		lastName string,
		phoneNumber string,
	) error

	UpdateUser(
		ctx context.Context, 
		id uint,
		username string,
		role string,
		firstName string,
		lastName string,
		phoneNumber string,
	) error

	DisableUser(ctx context.Context, ID uint) error
	EnableUser(ctx context.Context, ID uint) error
	GetUserByID(ctx context.Context, ID uint) (*m.User, error)
	ListUsers(ctx context.Context, ) ([]m.User, error)
}

type userServiceImpl struct {
	repo r.UserRepository
}

func NewUserService(repo r.UserRepository) UserService{
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) usernameValidation(
	ctx context.Context, 
	id uint, 
	username string,
) error {
	if len(username) < 4 {
		return e.NewBadRequest("username must be more than 4 characters long")
	}

	existingUser, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return err
	}

	if existingUser != nil && existingUser.ID != id {
		return e.NewConflict("username already exists")
	}

	return nil
}

func (s *userServiceImpl) RegisterUser(
	ctx context.Context, 
	username string,
	password string,
	role string,
	firstName string,
	lastName string,
	phoneNumber string,
) error {
	if err := s.usernameValidation(ctx, 0, username); err != nil{
		return err
	}

	if len(password) < 8 {
		return e.NewBadRequest("password must be at least 8 characters long")
	}
	
	password, err := hashPassword(password)
	if err != nil {
		return e.NewInternalServerError("password encryption failed")
	}

	var validRole m.UserRole
	switch role {
	case string(m.RoleCashier): validRole = m.RoleCashier
	case string(m.RoleScanner): validRole = m.RoleScanner
	default: return e.NewBadRequest("invalid user role")
	}

	phoneNumber, err = generateFormattedThaiPhone(phoneNumber)
	if err != nil {
		return e.NewBadRequest("invalid phone number")
	}

	user := &m.User{
		Username: username,
    	Password: password,
 		Role: validRole,
    	FirstName: firstName,
    	LastName: lastName,
    	PhoneNumber: phoneNumber,
    	IsActive: true,
	}

	return s.repo.Create(ctx, user)
}

func (s *userServiceImpl) UpdateUser(
	ctx context.Context, 
	id uint,
	username string,
	role string,
	firstName string,
	lastName string,
	phoneNumber string,
) error {
	if err := s.usernameValidation(ctx, id, username); err != nil{
		return err
	}

	var validRole m.UserRole
	switch role {
	case string(m.RoleCashier): validRole = m.RoleCashier
	case string(m.RoleScanner): validRole = m.RoleScanner
	default: return e.NewBadRequest("invalid user role")
	}

	phoneNumber, err := generateFormattedThaiPhone(phoneNumber)
	if err != nil {
		return e.NewBadRequest("invalid phone number")
	}

	user := &m.User{
		Username: username,
 		Role: validRole,
    	FirstName: firstName,
    	LastName: lastName,
    	PhoneNumber: phoneNumber,
	}

	return s.repo.Update(ctx, user)
}

func (s *userServiceImpl) DisableUser(ctx context.Context, id uint) error {
	return s.repo.SetActive(ctx, id, false)
}

func (s *userServiceImpl) EnableUser(ctx context.Context, id uint) error {
	return s.repo.SetActive(ctx, id, true)
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, id uint) (*m.User, error) {
	user, err := s.repo.GetByID(ctx,id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, e.NewNotFound("user not found")
	}
	return user, nil
}

func (s *userServiceImpl) ListUsers(ctx context.Context, ) ([]m.User, error) {
	return s.repo.List(ctx)
}