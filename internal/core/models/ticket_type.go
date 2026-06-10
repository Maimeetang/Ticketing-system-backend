package domains

import "time"

type TicketType struct {
	ID          uint       `gorm:"primaryKey"`
	Name        string     `gorm:"unique;not null"`
	Price       int64      `gorm:"not null"`
	Description string
	IsActive    bool       `gorm:"not null;default:true;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}