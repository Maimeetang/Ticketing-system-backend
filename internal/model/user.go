package model

import "time"

type User struct {
	ID                 uint
	role               string
	FirstName          string
	LastName           string
	PhoneNumber        uint
	ReservePhoneNumber uint
	CreatedAt          time.Time
	UpdatedAt 		   time.Time
}