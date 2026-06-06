package orm

import (
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
	return handleError(err)
}

func (r *GormUserRepository) UpdateUser(u *domain.User) error {
	result := r.db.Model(&domain.User{}).Updates(u)

	if result.Error != nil {
		return handleError(result.Error)
	}

	if result.RowsAffected == 0 {
		return apperror.NewNotFound("user")
	}

	return nil
}

func (r *GormUserRepository) SetActiveUser(id uint, active bool) error{
	result := r.db.
		Model(&domain.User{}).
		Where("id = ?", id).
		Update("is_active", active)

	if result.RowsAffected == 0 {
		return apperror.NewNotFound("user")
	}

	return handleError(result.Error)
}

func (r *GormUserRepository) GetUserByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, handleError(err) 
	}
	return &user, nil
}

func (r *GormUserRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, handleError(err)
	}
	return &user, nil
}

func (r *GormUserRepository) ListUsers() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, handleError(err)
	}
	return users, nil
}