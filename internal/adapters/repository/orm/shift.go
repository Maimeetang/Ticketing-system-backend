package orm

import (
	"errors"
	"strings"
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

// mapShiftDBError intercepts structural storage failure states
func mapShiftDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.NewNotFound("shift records not found")
	}

	// Capture foreign key constraint validation failure (e.g., clocking in a non-existent user)
	if isForeignKeyViolation(err) {
		return apperror.NewBadRequest("associated user ID does not exist")
	}

	return apperror.NewInternalServerError("database processing fault: " + err.Error())
}

// Auxiliary helper to catch constraint validation across relational storage engines
func isForeignKeyViolation(err error) bool {
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "foreign key") || strings.Contains(errMsg, "constraint failed")
}

func (r *GormShiftRepository) Create(shift *domain.Shift) error {
	err := r.db.Create(shift).Error
	return mapShiftDBError(err)
}

func (r *GormShiftRepository) Update(shift *domain.Shift) error {
	// Updates using struct attributes matching primaryKey criteria
	result := r.db.Model(&domain.Shift{}).Where("id = ?", shift.ID).Updates(shift)
	if result.Error != nil {
		return mapShiftDBError(result.Error)
	}
	
	if result.RowsAffected == 0 {
		return apperror.NewNotFound("target shift record missing")
	}
	return nil
}

func (r *GormShiftRepository) FindActiveByUserID(userID uint) (*domain.Shift, error) {
	var shift domain.Shift
	
	// Query database to find a shift that belongs to the user and is still 'OPEN'
	err := r.db.Where("user_id = ? AND status = ?", userID, domain.ShiftOpen).First(&shift).Error
	if err != nil {
		return nil, mapShiftDBError(err)
	}
	
	return &shift, nil
}
