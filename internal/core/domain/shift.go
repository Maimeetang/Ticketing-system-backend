package domain

import "time"

type ShiftStatus string

const (
	ShiftOpen   ShiftStatus = "OPEN"
	ShiftClosed ShiftStatus = "CLOSED"
)

type Shift struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	UserID    uint        `gorm:"not null;index:idx_user_status;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user_id"`
	StartAt   time.Time   `gorm:"not null" json:"start_at"`
	EndAt     *time.Time  `gorm:"default:null" json:"end_at"`
	Status    ShiftStatus `gorm:"type:varchar(20);not null;index:idx_user_status;default:'OPEN'" json:"status"`
}
