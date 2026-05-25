package domain

import "time"

type TicketInfo struct {
	ID           uint       `gorm:"primaryKey"`
	TicketID  	 uint       `gorm:"not null;index"`
	TicketTypeID uint		`gorm:"not null;index"`
	PricePerUnit float64    `gorm:"type:decimal(10,2);not null"`
	Quantity     int        `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}