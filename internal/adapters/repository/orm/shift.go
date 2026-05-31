package orm

import (
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"

	"gorm.io/gorm"
)

type GormShiftRepository struct {
	db *gorm.DB
}

func NewGormShiftRepository(db *gorm.DB) port.ShiftRepository {
	return &GormShiftRepository{db: db}
}

func (r *GormShiftRepository) CreateShift(shift *domain.Shift) error {
	err := r.db.Create(shift).Error
	return handleError(err)
}

func (r *GormShiftRepository) UpdateShift(shift *domain.Shift) error {
	result := r.db.Model(&domain.Shift{}).Where("id = ?", shift.ID).Updates(shift)
	if result.Error != nil {
		return handleError(result.Error)
	}
	
	if result.RowsAffected == 0 {
		return apperror.NewNotFound("target shift record missing")
	}
	return nil
}

func (r *GormShiftRepository) GetActiveByUserID(userID uint) (*domain.Shift, error) {
	var shift domain.Shift
	
	err := r.db.Where("user_id = ? AND status = ?", userID, domain.ShiftOpen).First(&shift).Error
	if err != nil {
		return nil, handleError(err)
	}

	return &shift, nil
}