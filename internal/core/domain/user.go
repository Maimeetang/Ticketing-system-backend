package domain

import "time"

type User struct {
	ID                 uint      `json:"id"`
	Username           string    `json:"username"`
	Password           string    `json:"password"`
	Role      		   string    `json:"role"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	PhoneNumber        string    `json:"phone_number"`
	ReservePhoneNumber string    `json:"reserve_phone_number"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}