package port

import "ticketing-system/internal/core/domain"

type ShiftService interface {
	OpenShift(userID uint) (*domain.Shift, error)
	CloseShift(id uint) error
	GetCurrentShift(userID uint) (*domain.Shift, error)
	GetShiftByID(id uint) (*domain.Shift, error)
	ListShifts(domain.ShiftFilter) ([]domain.Shift, error)
}

type ShiftRepository interface {
	Create(shift *domain.Shift) error
	Update(shift *domain.Shift) error
	GetByID(id uint) (*domain.Shift, error)
	GetCurrentByUserID(userID uint) (*domain.Shift, error)
	List(domain.ShiftFilter) ([]domain.Shift, error)
}
