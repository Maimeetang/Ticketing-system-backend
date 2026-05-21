package http

import (
	"ticketing-system/internal/core/domain"
	"time"
)

// userResponse represents a user response body
type userResponse struct {
	ID                 uint   `json:"id"`
	Username           string `json:"username"`
	Role               string `json:"role"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	PhoneNumber        string `json:"phone_number"`
	ReservePhoneNumber string `json:"reserve_phone_number"`
}

// newUserResponse is a helper function to create a user body for handling user data
func newUserResponse(user *domain.User) userResponse {
	return userResponse{
		ID: 	  			user.ID,
		Username: 			user.Username,
		Role: 				string(user.Role),
		FirstName: 			user.FirstName,
		LastName: 			user.LastName,
		PhoneNumber: 		user.PhoneNumber,
		ReservePhoneNumber: user.ReservePhoneNumber,
	}
}

// shiftResponse represents a shift response body
type shiftResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	StartAt   time.Time `json:"start_at"`
	Status    string    `json:"status"`
}

// newShiftResponse is a helper function to create a response body for handling shift data
func newShiftResponse(shift *domain.Shift) shiftResponse {
	return shiftResponse{
		ID: 	 shift.ID,
		UserID:  shift.UserID,
		StartAt: shift.StartAt,
		Status:  string(shift.Status),
	}
}
