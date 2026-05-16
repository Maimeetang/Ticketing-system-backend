package model

import "time"

type User struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	Username           string    `gorm:"unique" json:"username"`
	Password           string    `json:"-"`
	Role      		   string    `gorm:"not null" json:"role"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	PhoneNumber        string    `gorm:"unique" json:"phone_number"`
	ReservePhoneNumber string    `json:"reserve_phone_number"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}