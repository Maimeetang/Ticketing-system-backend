package orm

import (
	"errors"
	"strings"
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"

	"gorm.io/gorm"
)

type GormUserRepository struct {
  	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) port.UserRepository {
  	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) CreateUser(u *domain.User) error {
	err := r.db.Create(u).Error
	return handleUserError(err)
}

func (r *GormUserRepository) UpdateUser(u *domain.User) error {
	result := r.db.Model(&domain.User{}).Where("id = ?", u.ID).Updates(u)

	if result.Error != nil {
		return handleUserError(result.Error)
	}

	if result.RowsAffected == 0 {
		return apperror.NewNotFound("user")
	}

	return nil
}

func (r *GormUserRepository) DeleteUser(id uint) error {
	result := r.db.Delete(&domain.User{}, id)

	if result.Error != nil {
		return handleUserError(result.Error)
	}
	if result.RowsAffected == 0 {
		return apperror.NewNotFound("user")
	}
	return nil
}

func (r *GormUserRepository) GetUserByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, handleUserError(err) 
	}
	return &user, nil
}

func (r *GormUserRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, handleUserError(err)
	}
	return &user, nil
}

func (r *GormUserRepository) ListUsers() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, handleUserError(err)
	}
	return users, nil
}

func handleUserError(err error) error {
	if err == nil {
		return nil
	}
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.NewNotFound(err.Error())
	}

	if strings.Contains(strings.ToUpper(err.Error()), "UNIQUE") {
		if strings.Contains(strings.ToLower(err.Error()), "username") {
			return apperror.NewConflict("username already exists")
		}
		return apperror.NewConflict("data conflict")
	}

	return apperror.NewInternalServerError("database error: " + err.Error())
}