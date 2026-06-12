package db

import (
	"context"
	e "ticketing-system/internal/core/error"
	m "ticketing-system/internal/core/models"
	r "ticketing-system/internal/core/repositories"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormShiftRepository struct {
	db *gorm.DB
}

func NewGormShiftRepository(db *gorm.DB) r.ShiftRepository {
	return &GormShiftRepository{db: db}
}

func (r *GormShiftRepository) Create(ctx context.Context, shift *m.Shift) error {
	db := bind(ctx, r.db)
	return handleError(db.Create(shift).Error)
}

func (r *GormShiftRepository) Update(ctx context.Context, shift *m.Shift) error {
	db := bind(ctx, r.db)

	result := db.
			Where("id = ?", shift.ID).
			Model(&m.Shift{}).
			Updates(shift)
	if result.Error != nil {
		return handleError(result.Error)
	}
	
	if result.RowsAffected == 0 {
		return e.NewNotFound("ไม่พบกะทำงาน")
	}
	return nil
}

func (r *GormShiftRepository) GetByID(ctx context.Context, id uint) (*m.Shift, error) {
	var shift m.Shift
	db := bind(ctx, r.db)

	err := db.
		Preload("User").
		First(&shift, id).
		Error

	if err != nil {
		return nil, handleError(err)
	}

	return &shift, nil
}

func (r *GormShiftRepository) GetByIDForUpdate(
	ctx context.Context, 
	id uint,
) (*m.Shift, error) {
	var shift m.Shift
	db := bind(ctx, r.db)

	err := db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&shift, id).Error

	return &shift, handleError(err)
}

func (r *GormShiftRepository) GetCurrentByUserID(
	ctx context.Context,
	userID uint,
) (*m.Shift, error) {
	var shift m.Shift
	db := bind(ctx, r.db)

	err := db.
		Preload("User").
		Where(
			"user_id = ? AND status = ?",
			userID,
			m.ShiftOpen,
		).
		First(&shift).
		Error

	if err != nil {
		return nil, handleError(err)
	}

	return &shift, nil
}

func (r *GormShiftRepository) GetCurrentByUserIDForUpdate(
	ctx context.Context,
	userID uint,
) (*m.Shift, error) {
	var shift m.Shift
	db := bind(ctx, r.db)

	err := db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(
			"user_id = ? AND status = ?",
			userID,
			m.ShiftOpen,
		).
		First(&shift).Error

	return &shift, handleError(err)
}

func (r *GormShiftRepository) GetByDate(
	ctx context.Context, 
	targetDate time.Time,
) ([]m.Shift, error) {
	var shifts []m.Shift

	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Second)

	err := bind(ctx, r.db).
		Preload("User").
		Where("open_at BETWEEN ? AND ?", startOfDay, endOfDay).
		Find(&shifts).
		Error

	if err != nil {
		return nil, handleError(err)
	}

	return shifts, nil
}

func (r *GormShiftRepository) CalculateSummary(ctx context.Context, shift *m.Shift) error {
	var total int64
	var cancelled int64
	var revenue int64

	db := bind(ctx, r.db)

	if err := db.Model(&m.Ticket{}).Where("shift_id = ?", shift.ID).Count(&total).Error; err != nil {
		return handleError(err)
	}

	if err := db.Model(&m.Ticket{}).Where("shift_id = ? AND status = ?", shift.ID, m.TicketCancelled).Count(&cancelled).Error; err != nil {
		return handleError(err)
	}

	if err := db.Model(&m.Ticket{}).Where("shift_id = ? AND status <> ?", shift.ID, m.TicketCancelled).Select("COALESCE(SUM(total_price), 0)").Scan(&revenue).Error; err != nil {
		return handleError(err)
	}

	shift.TotalTickets = uint(total)
	shift.CancelledTickets = uint(cancelled)
	shift.TotalRevenue = revenue

	return nil
}
