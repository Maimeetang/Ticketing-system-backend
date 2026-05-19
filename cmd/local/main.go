package main

import (
	"log"
	"os"
	"ticketing-system/internal/adapters/handler/http"
	"ticketing-system/internal/adapters/repository/orm"
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/config"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/service"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Load configurations from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("env file missing; system will fall back to environment variables.")
	}

	cfg := config.LoadAuthConfig()
	appMode := os.Getenv("APP_MODE")
	if appMode == "" {
		appMode = "offline" // Default mode
	}

	log.Printf("Initializing Ticketing System in [%s] mode...", appMode)

	// Connect to Local SQLite Database File
	// It will create 'local_tickets.db' file automatically in root folder if it doesn't exist.
	db, err := gorm.Open(sqlite.Open("./data/local_tickets.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database initialization aborted (SQLite): %v", err)
	}
	log.Println("Local SQLite Database connected successfully.")

	// Auto-Migrate Schemas directly through GORM
	// It will scan domain.User structural constraints and map to SQLite table architecture.
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("Database schema migration failed: %v", err)
	}
	log.Println("Database schemas auto-migrated successfully.")

	// 4. Orchestrate Dependency Injection (Wiring layers together)
	userRepo := orm.NewGormUserRepository(db) // SQLite uses the same GormUserRepository since it is an ORM abstraction layer
	
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, cfg)

	userHandler := http.NewUserHandler(userService)
	authHandler := http.NewAuthHandler(authService, cfg) // Injected cfg to allow dynamic Cookie Secure tuning

	// 5. Initialize Fiber Web Application and Centralized Error Mapping Middleware
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			statusCode := fiber.StatusInternalServerError
			message := "Internal server fault occurred"

			// Capture internal decoupled application domain errors
			if appErr, ok := err.(apperror.AppError); ok {
				statusCode = appErr.StatusCode()
				message = appErr.Message()
			} else if fiberErr, ok := err.(*fiber.Error); ok {
				statusCode = fiberErr.Code
				message = fiberErr.Message
			}

			return c.Status(statusCode).JSON(fiber.Map{
				"status":  "error",
				"message": message,
			})
		},
	})

	// 6. Register Routing Schema Architecture
	api := app.Group("/api")
	
	// Public Endpoints
	api.Post("/auth/login", authHandler.Login)
	api.Post("/auth/logout", authHandler.Logout)

	// Protected Endpoints (Guarded by Cookie Verification Gate)
	userRoutes := api.Group("/users", http.JWTMiddleware(cfg))
	
	userRoutes.Post("/", userHandler.AddUser)
	userRoutes.Put("/:id", userHandler.EditUser)
	userRoutes.Delete("/:id", userHandler.DeleteUser)
	userRoutes.Get("/:id", userHandler.FindUserByID)
	userRoutes.Get("/", userHandler.ListUsers)

	// 7. Fire up the local web engine
	log.Fatal(app.Listen(":8000"))
}
