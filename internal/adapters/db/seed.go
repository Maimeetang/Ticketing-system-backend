package db

import (
	"log"
	m "ticketing-system/internal/core/models"
	s "ticketing-system/internal/core/services"

	"gorm.io/gorm"
)

func SeedAdminAccount(
	gormDB *gorm.DB,
	uName string,
	pass string,
	fName string,
	lName string,
	pNumber string,
) {
	var count int64
	err := gormDB.
		Model(&m.Shift{}).
		Table("users").
		Where("role = ?", "ADMIN").
		Count(&count).Error

	if err != nil {
		log.Printf("Failed to check existing admin account: %v", err)
		return
	}

	if count > 0 {
		log.Println("Admin account already exists. Skipping seeding.")
		return
	}

	hashedPassword, err := s.HashPassword(pass)
	if err != nil {
		log.Printf("Failed to hash seed admin password: %v", err)
		return
	}

	adminUser := &m.User{
    	Username: uName,
    	Password: hashedPassword,
    	Role: m.RoleAdmin, 
    	FirstName: fName,
    	LastName : lName, 
    	PhoneNumber: pNumber,
		IsActive: true,
	}

	err = gormDB.Create(&adminUser).Error
	if err != nil {
		log.Printf("Failed to seed default admin account: %v", err)
		return
	}

	log.Println("Successfully seeded default admin account")
}