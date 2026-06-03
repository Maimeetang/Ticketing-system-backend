package domain

import "time"

type User struct {
	ID                 uint      `gorm:"primaryKey"`
	Username           string    `gorm:"unique;not null"`
	Password           string    `gorm:"not null"`
	Role      		   UserRole  `gorm:"not null"`
	FirstName          string
	LastName           string
	PhoneNumber        string
	ReservePhoneNumber string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}