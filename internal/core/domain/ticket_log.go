package domain

import "time"

type TicketLog struct {
	ID          uint          `gorm:"primaryKey"`
	UserID		uint		  `gorm:"not null"`
	TicketID    uint          `gorm:"not null;index"`
	FromStatus  *TicketStatus `gorm:"type:varchar(20);"`
	ToStatus    TicketStatus  `gorm:"type:varchar(20);not null"`
	Remarks     string        `gorm:"type:varchar(255)"`
	CreatedAt   time.Time     `gorm:"not null;index"`
}