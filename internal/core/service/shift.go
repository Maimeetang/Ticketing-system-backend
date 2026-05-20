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
func (s *shiftServiceImpl) ClockIn(userID uint) (*domain.Shift, error) {
	activeShift, err := s.repo.FindActiveByUserID(userID)
	if err == nil && activeShift != nil {
		return nil, apperror.NewConflict("user already has an active working shift session open")
	}

	newShift := &domain.Shift{
		UserID:    userID,
		StartAt:   time.Now(),
		EndAt:     nil,
		Status:    domain.ShiftOpen,
	}

	if err := s.repo.Create(newShift); err != nil {
		return nil, err
	}

	return newShift, nil
}

func (s *shiftServiceImpl) ClockOut(userID uint) error {
	activeShift, err := s.repo.FindActiveByUserID(userID)
	if err != nil {
		return apperror.NewNotFound("no active shift session found for this user to clock out")
	}

	now := time.Now()
	activeShift.EndAt = &now
	activeShift.Status = domain.ShiftClosed

	return s.repo.Update(activeShift)
}

func (s *shiftServiceImpl) GetActiveShift(userID uint) (*domain.Shift, error) {
	activeShift, err := s.repo.FindActiveByUserID(userID)
	if err != nil {
		return nil, apperror.NewNotFound("no active shift running for this user")
	}
	return activeShift, nil
}