package domains

import "time"

type UserRole string

const (
	RoleAdmin   UserRole = "ADMIN"
	RoleManager UserRole = "MANAGER"
	RoleCashier UserRole = "CASHIER"
	RoleScanner UserRole = "SCANNER"
)

type User struct {
	ID uint `gorm:"primaryKey"`

	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`

	Role UserRole `gorm:"not null"`

	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`

	PhoneNumber 	   string `gorm:"not null"`

	IsActive  bool `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}