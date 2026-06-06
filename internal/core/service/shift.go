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

func (s *shiftServiceImpl) OpenShift(userID uint) (*domain.Shift, error) {
	currentShift, err := s.repo.GetCurrentByUserID(userID)
	if err != nil {
		return nil, err
	}

	if currentShift != nil {
		return nil, apperror.NewConflict("คุณมีกะการทำงานที่เปิดอยู่แล้ว")
	}

	shift := &domain.Shift{
		UserID: userID,
		OpenAt: time.Now(),
		Status: domain.ShiftOpen,
	}

	if err := s.repo.Create(shift); err != nil {
		return nil, err
	}

	return shift, nil
}

func (s *shiftServiceImpl) CloseShift(id uint)  error {
	shift, err := s.repo.GetByID(id)

	if err != nil {
		return err
	}

	if shift == nil {
		return apperror.NewNotFound("ไม่พบกะการทำงานที่กำลังเปิดอยู่")
	}

	if shift.Status == domain.ShiftClosed {
		return apperror.NewConflict("กะทำงานนี้ถูกปิดอยู่แล้ว")
	}

	now := time.Now()

	shift.CloseAt = &now
	shift.Status = domain.ShiftClosed

	return s.repo.Update(shift)
}

func (s *shiftServiceImpl) GetCurrentShift(userID uint) (*domain.Shift, error) {
	shift, err := s.repo.GetCurrentByUserID(userID)
	if err != nil {
		return nil, err
	}

	if shift == nil {
		return nil, apperror.NewNotFound("ไม่พบกะการทำงานที่กำลังเปิดอยู่")
	}

	return shift, nil
}

func (s *shiftServiceImpl) GetShiftByID(id uint) (*domain.Shift, error) {
	shift, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if shift == nil {
		return nil, apperror.NewNotFound("ไม่พบกะการทำงาน")
	}

	return shift, nil
}

func (s *shiftServiceImpl) ListShifts(filter domain.ShiftFilter) ([]domain.Shift, error) {
	return s.repo.List(filter)
}