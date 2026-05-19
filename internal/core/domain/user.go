package domain

import "time"

type UserRole string

const (
	RoleAdmin 	UserRole = "admin"
	RoleManager UserRole = "manager"
	RoleCashier UserRole = "cashier"
	RoleScanner UserRole = "scanner"
)

type User struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	Username           string    `gorm:"unique;not null" json:"username"`
	Password           string    `gorm:"not null" json:"password"`
	Role      		   UserRole  `gorm:"not null" json:"role"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	PhoneNumber        string    `json:"phone_number"`
	ReservePhoneNumber string    `json:"reserve_phone_number"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}