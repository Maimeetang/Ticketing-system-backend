package services

import (
	"context"
	e "ticketing-system/internal/core/error"
	m "ticketing-system/internal/core/models"
	r "ticketing-system/internal/core/repositories"
	"time"
)

type ShiftService interface {
	OpenShift(ctx context.Context, userID uint) (*m.Shift, error)
	CloseShift(ctx context.Context, id uint) error
	GetCurrentShift(ctx context.Context, userID uint) (*m.Shift, error)
	GetShiftByID(ctx context.Context, id uint) (*m.Shift, error)
	GetShiftByDate(ctx context.Context, date string) ([]m.Shift, error)
}

type shiftServiceImpl struct {
	transactor r.Transactor
	shiftRepo r.ShiftRepository
	userRepo r.UserRepository
}

func NewShiftService(
	transactor r.Transactor, 
	shiftRepo r.ShiftRepository,
	userRepo r.UserRepository,
) ShiftService {
	return &shiftServiceImpl{
		transactor: transactor,
		shiftRepo: shiftRepo,
		userRepo: userRepo,
	}
}

func (s *shiftServiceImpl) OpenShift(ctx context.Context, userID uint) (*m.Shift, error) {
	var shift *m.Shift

	err := s.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		_, err := s.userRepo.GetByIDForUpdate(txCtx, userID)
		if err != nil {
			return err
		}

		currentShift, err := s.shiftRepo.GetCurrentByUserID(txCtx, userID)
		if err != nil {
			return err
		}

		if currentShift != nil {
			return e.NewConflict("shift is already open")
		}

		now := time.Now()
		shift = &m.Shift{
			UserID: userID,
			OpenAt: now,
			Status: m.ShiftOpen,
		}

		return s.shiftRepo.Create(txCtx, shift)
	})

	if err != nil {
		return nil, err
	}
	return shift, nil
}

func (s *shiftServiceImpl) CloseShift(ctx context.Context, id uint) error {
	return s.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		shift, err := s.shiftRepo.GetByIDForUpdate(txCtx, id)
		if err != nil {
			return err
		}

		if shift == nil {
			return e.NewNotFound("shift not found")
		}

		if shift.Status == m.ShiftClosed {
			return e.NewConflict("shift has already been closed")
		}

		if err := s.shiftRepo.CalculateSummary(txCtx, shift); err != nil {
			return err
		}

		now := time.Now()
		shift.CloseAt = &now
		shift.Status = m.ShiftClosed

		if err := s.shiftRepo.Update(txCtx, shift); err != nil {
			return err
		}

		return nil
	})
}

func (s *shiftServiceImpl) GetCurrentShift(
	ctx context.Context, userID uint,
) (*m.Shift, error) {
	shift, err := s.shiftRepo.GetCurrentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if shift == nil {
		return nil, e.NewNotFound("shift not found")
	}

	return shift, nil
}

func (s *shiftServiceImpl) GetShiftByID(
	ctx context.Context, id uint,
) (*m.Shift, error) {
	shift, err := s.shiftRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if shift == nil {
		return nil, e.NewNotFound("shift not found")
	}

	return shift, nil
}

func (s *shiftServiceImpl) GetShiftByDate(
	ctx context.Context, 
	dateStr string,
) ([]m.Shift, error) {
	targetDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, e.NewBadRequest("invalid date format, please use YYYY-MM-DD")
	}

	shifts, err := s.shiftRepo.GetByDate(ctx, targetDate)
	if err != nil {
		return nil, err
	}

	return shifts, nil
}