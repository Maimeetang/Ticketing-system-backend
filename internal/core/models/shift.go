package domains

import "time"

type ShiftStatus string

const (
	ShiftOpen   ShiftStatus = "OPEN"
	ShiftClosed ShiftStatus = "CLOSED"
)

type Shift struct {
	ID		uint `gorm:"primaryKey"`
	UserID 	uint `gorm:"not null;index:idx_user_status"`

	OpenAt	time.Time `gorm:"not null"`
	CloseAt *time.Time

	Status ShiftStatus `gorm:"type:varchar(20);not null;index:idx_user_status;default:'OPEN'"`

	TotalTickets 	 uint `gorm:"not null;default:0"`
	CancelledTickets uint `gorm:"not null;default:0"`

	TotalRevenue int64 `gorm:"not null;default:0"`

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type ShiftFilter struct {
	UserID *uint
	Status *ShiftStatus
	StartDate time.Time
}
