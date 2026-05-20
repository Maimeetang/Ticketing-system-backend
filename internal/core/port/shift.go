package port

import "ticketing-system/internal/core/domain"

// Primary port
type ShiftService interface {
	ClockIn(userID uint) (*domain.Shift, error)
	ClockOut(userID uint) error
	GetActiveShift(userID uint) (*domain.Shift, error)
}

// Secondary port
type ShiftRepository interface {
	Create(shift *domain.Shift) error
	Update(shift *domain.Shift) error
	FindActiveByUserID(userID uint) (*domain.Shift, error)
}
