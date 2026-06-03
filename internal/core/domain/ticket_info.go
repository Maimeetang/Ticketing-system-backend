package domain

import "time"

type TicketInfo struct {
	ID           uint       `gorm:"primaryKey"`
	TicketID  	 uint       `gorm:"uniqueIndex"`
	TicketType   string		`gorm:"not null"`
	PricePerUnit float64    `gorm:"type:decimal(10,2);not null"`
	Quantity     int        `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}