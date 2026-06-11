package db

import (
	"context"
	e "ticketing-system/internal/core/error"
	m "ticketing-system/internal/core/models"
	r "ticketing-system/internal/core/repositories"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormUserRepository struct {
  	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) r.UserRepository {
  	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(ctx context.Context, u *m.User) error {
	db := bind(ctx, r.db)

	err := db.Create(u).Error
	return handleError(err)
}

func (r *GormUserRepository) Update(ctx context.Context, u *m.User) error {
	db := bind(ctx, r.db)

	result := db.Updates(u)

	if result.Error != nil {
		return handleError(result.Error)
	}

	if result.RowsAffected == 0 {
		return e.NewNotFound("user not found")
	}

	return nil
}

func (r *GormUserRepository) SetActive(ctx context.Context, id uint, active bool) error{
	db := bind(ctx, r.db)

	result := db.
		Model(&m.User{}).
		Where("id = ?", id).
		Update("is_active", active)

	if result.RowsAffected == 0 {
		return e.NewNotFound("user not found")
	}

	return handleError(result.Error)
}

func (r *GormUserRepository) GetByID(ctx context.Context, id uint) (*m.User, error) {
	var user m.User
	db := bind(ctx, r.db)

	err := db.First(&user, id).Error
	if err != nil {
		return nil, handleError(err) 
	}
	return &user, nil
}

func (r *GormUserRepository) GetByIDForUpdate(ctx context.Context, id uint) (*m.User, error) {
	var user m.User
	db := bind(ctx, r.db)

	err := db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&user, id).Error

	return &user, handleError(err)
}

func (r *GormUserRepository) GetByUsername(
	ctx context.Context, 
	username string,
) (*m.User, error) {
	var user m.User
	db := bind(ctx, r.db)

	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, handleError(err)
	}
	return &user, nil
}

func (r *GormUserRepository) List(ctx context.Context, ) ([]m.User, error) {
	var users []m.User
	db := bind(ctx, r.db)

	err := db.Find(&users).Error
	if err != nil {
		return nil, handleError(err)
	}
	return users, nil
}