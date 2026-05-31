package service

import (
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"
	"time"
)

type shiftServiceImpl struct {
	repo port.ShiftRepository
}

func NewShiftService(repo port.ShiftRepository) port.ShiftService {
	return &shiftServiceImpl{repo: repo}
}

// Business Core logic
func (s *shiftServiceImpl) ClockIn(shift *domain.Shift) (*domain.Shift, error) {
	activeShift, err := s.repo.GetActiveByUserID(shift.UserID)
	if err != nil {
		return nil, err
	}

	if activeShift != nil {
		return nil, apperror.NewConflict("คุณมีกะการทำงานที่เปิดอยู่แล้ว")
	}

	if err := s.repo.CreateShift(shift); err != nil {
		return nil, err
	}

	return shift, nil
}

func (s *shiftServiceImpl) ClockOut(userID uint) error {
	shift, err := s.repo.GetActiveByUserID(userID)

	if err != nil {
		return err
	}

	if shift == nil {
		return apperror.NewNotFound("ไม่พบกะการทำงานที่กำลังเปิดอยู่")
	}

	now := time.Now()
	shift.EndAt = &now
	shift.Status = domain.ShiftClosed

	return s.repo.UpdateShift(shift)
}

func (s *shiftServiceImpl) GetUserActiveShift(userID uint) (*domain.Shift, error) {
	shift, err := s.repo.GetActiveByUserID(userID)

	if err != nil {
		return nil, err
	}

	return shift, nil
}