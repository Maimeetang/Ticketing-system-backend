package dto

import "time"

type ShiftResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	StartAt   time.Time `json:"start_at"`
	Status    string    `json:"status"`
}
