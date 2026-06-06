package orm

import (
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"
	"time"

	"gorm.io/gorm"
)

type GormShiftRepository struct {
	db *gorm.DB
}

func NewGormShiftRepository(db *gorm.DB) port.ShiftRepository {
	return &GormShiftRepository{db: db}
}

func (r *GormShiftRepository) Create(shift *domain.Shift) error {
	return handleError(r.db.Create(shift).Error)
}

func (r *GormShiftRepository) Update(shift *domain.Shift) error {
	result := r.db.
			Where("id = ?", shift.ID).
			Model(&domain.Shift{}).
			Updates(shift)
	if result.Error != nil {
		return handleError(result.Error)
	}
	
	if result.RowsAffected == 0 {
		return apperror.NewNotFound("ไม่พบกะทำงาน")
	}
	return nil
}

func (r *GormShiftRepository) GetByID(id uint) (*domain.Shift, error) {
	var shift domain.Shift

	err := r.db.
		Preload("User").
		First(&shift, id).
		Error

	if err != nil {
		return nil, handleError(err)
	}

	return &shift, nil
}

func (r *GormShiftRepository) GetCurrentByUserID(userID uint,) (*domain.Shift, error) {
	var shift domain.Shift

	err := r.db.
		Preload("User").
		Where(
			"user_id = ? AND status = ?",
			userID,
			domain.ShiftOpen,
		).
		First(&shift).
		Error

	if err != nil {
		return nil, handleError(err)
	}

	return &shift, nil
}

func (r *GormShiftRepository) List(filter domain.ShiftFilter,) ([]domain.Shift, error) {
	var shifts []domain.Shift

	query := r.db.
		Preload("User")

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	start := time.Date(
			filter.StartDate.Year(),
			filter.StartDate.Month(),
			filter.StartDate.Day(),
			0, 0, 0, 0,
			filter.StartDate.Location(),
		)

	end := start.AddDate(0, 0, 1)

	query = query.Where(
		"open_at >= ? AND open_at < ?",
		start,
		end,
	)

	err := query.
		Order("open_at DESC").
		Find(&shifts).
		Error

	if err != nil {
		return nil, handleError(err)
	}

	return shifts, nil
}