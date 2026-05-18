package orm

import (
	"errors"
	"strings"
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/model"
	"ticketing-system/internal/core/repository"

	"gorm.io/gorm"
)

type GormUserRepository struct {
  	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) repository.UserRepository {
  	return &GormUserRepository{db: db}
}

func mapDBError(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "UNIQUE") {
		if strings.Contains(err.Error(), "username") {
			return apperror.NewConflict("username")
		}
		return apperror.NewConflict("data")
	}
	return err
}

func (r *GormUserRepository) Add(user model.User) error {
	err := r.db.Create(&user).Error
	return mapDBError(err)
}

func (r *GormUserRepository) Edit(user model.User) error {
	err := r.db.Save(&user).Error
	return mapDBError(err)
}

func (r *GormUserRepository) DeleteByID(id uint) error {
	result := r.db.Delete(&model.User{}, id)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperror.NewNotFound("user")
	}
	return nil
}

func (r *GormUserRepository) FindByID(id uint) (model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, apperror.NewNotFound("user")
	}
	return user, err
}

func (r *GormUserRepository) FindByUsername(username string) (model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, apperror.NewNotFound("user")
	}
	return user, err
}