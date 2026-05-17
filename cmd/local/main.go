package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ticketing-system/internal/adapters/handler"
	"ticketing-system/internal/adapters/orm"
	"ticketing-system/internal/core/model"
	"ticketing-system/internal/core/service"
)

func main() {
	app := fiber.New()

	// Initialize the database connection
	db, err := gorm.Open(sqlite.Open("./data/user.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrate the schema
  	db.AutoMigrate(&model.User{})

	userRepo := orm.NewGormUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	app.Post("/users", userHandler.AddUser)

	app.Listen(":8000")
}