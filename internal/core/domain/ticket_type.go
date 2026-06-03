package domain

import "time"

type TicketType struct {
	ID          uint       `gorm:"primaryKey"`
	Name        string     `gorm:"unique;not null"`
	Price       float64    `gorm:"type:decimal(10,2);not null"`
	Description string
	IsActive    bool       `gorm:"not null;default:true;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}