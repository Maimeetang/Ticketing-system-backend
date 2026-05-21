package port

import "ticketing-system/internal/core/domain"

// Primary port
type ShiftService interface {
	// ClockIn creates a new work shift for an employee.
	ClockIn(shift *domain.Shift) (*domain.Shift, error)
	// ClockOut ends the current active shift.
	ClockOut(userID uint) error
	// GetUserActiveShift retrieves the current active shift for a given user.
	GetUserActiveShift(userID uint) (*domain.Shift, error)
}

// Secondary port
type ShiftRepository interface {
	// CreateShift inserts a new shift record into the database.
	CreateShift(shift *domain.Shift) error
	// UpdateShift updates an existing shift.
	UpdateShift(shift *domain.Shift) error
	// GetActiveByUserID fetches the active shift for a specific user.
	GetActiveByUserID(userID uint) (*domain.Shift, error)
}
