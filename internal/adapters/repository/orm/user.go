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

func mapDBError(err error) error {
	if err == nil {
		return nil
	}
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.NewNotFound("user")
	}

	if strings.Contains(strings.ToUpper(err.Error()), "UNIQUE") {
		if strings.Contains(strings.ToLower(err.Error()), "username") {
			return apperror.NewConflict("username already exists")
		}
		return apperror.NewConflict("data conflict")
	}

	return apperror.NewInternalServerError("database error: " + err.Error())
}

func (r *GormUserRepository) Create(u domain.User) error {
	err := r.db.Create(&u).Error
	return mapDBError(err)
}

func (r *GormUserRepository) Update(u domain.User) error {
	result := r.db.Model(&domain.User{}).Where("id = ?", u.ID).Updates(&u)

	if result.Error != nil {
		return mapDBError(result.Error)
	}

	if result.RowsAffected == 0 {
		return apperror.NewNotFound("user")
	}

	return nil
}

func (r *GormUserRepository) DeleteByID(id uint) error {
	result := r.db.Delete(&domain.User{}, id)

	if result.Error != nil {
		return mapDBError(result.Error)
	}
	if result.RowsAffected == 0 {
		return apperror.NewNotFound("user")
	}
	return nil
}

func (r *GormUserRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, mapDBError(err) 
	}
	return &user, nil
}

func (r *GormUserRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, mapDBError(err)
	}
	return &user, nil
}

func (r *GormUserRepository) GetAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, mapDBError(err)
	}
	return users, nil
}