package main

import (
	"fmt"
	"log"
	"ticketing-system/config"
	seed "ticketing-system/internal/adapters/db"
	db "ticketing-system/internal/adapters/db/postgres"
	h "ticketing-system/internal/adapters/http/handlers"
	middleware "ticketing-system/internal/adapters/http/middleware"
	e "ticketing-system/internal/core/error"
	m "ticketing-system/internal/core/models"
	s "ticketing-system/internal/core/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB, cfg.PostgresPort,
	)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	log.Println("Database connected successfully")

	err = gormDB.AutoMigrate(&m.Shift{}, &m.Ticket{}, &m.TicketType{}, &m.User{}, &m.TicketLog{})
	if err != nil {
		log.Fatalf("AutoMigration failed: %v", err)
	}

	seed.SeedAdminAccount(
		gormDB, 
		cfg.DefaultAdminUsername, 
		cfg.DefaultAdminPassword,
		cfg.DefaultAdminFirstname,
		cfg.DefaultAdminLastname,
		cfg.DefaultAdminPhonenumber,
	)

	// ==========================================
	// INITIALIZE ADAPTERS (DATABASE LAYER)
	// ==========================================

	transactor := db.NewGormTransactor(gormDB)
	shiftRepo := db.NewGormShiftRepository(gormDB)
	ticketRepo := db.NewGormTicketRepository(gormDB)
	ticketTypeRepo := db.NewGormTicketTypeRepository(gormDB)
	userRepo := db.NewGormUserRepository(gormDB)
	ticketLogRepo := db.NewGormTicketLogRepository(gormDB)

	// ==========================================
	// INITIALIZE CORE SERVICES (BUSINESS LOGIC)
	// ==========================================

	authService := s.NewAuthService(userRepo)
	userService := s.NewUserService(userRepo)
	shiftService := s.NewShiftService(transactor, shiftRepo, userRepo)
	ticketService := s.NewTicketService(transactor, shiftRepo, ticketRepo, ticketLogRepo, ticketTypeRepo)
	ticketTypeService := s.NewTicketTypeService(ticketTypeRepo)

	// ==========================================
	// INITIALIZE HTTP HANDLERS (PRESENTATION LAYER)
	// ==========================================

	authHandler := h.NewAuthHandler(authService, cfg.AppEnv == "production", cfg.JwtSecret)
	userHandler := h.NewUserHandler(userService)
	shiftHandler := h.NewShiftHandler(shiftService)
	ticketHandler := h.NewTicketHandler(ticketService)
	ticketTypeHandler := h.NewTicketTypeHandler(ticketTypeService)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			statusCode := fiber.StatusInternalServerError
			message := "internal server processing fault occurred"

			if appErr, ok := err.(e.AppError); ok {
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

	api := app.Group("/api/v1")

	api.Post("/auth/login", authHandler.Login)
	api.Post("/auth/logout", authHandler.Logout)

	secured := api.Use(middleware.AuthRequired(cfg.JwtSecret))
	secured.Get("/auth/me", authHandler.GetProfile)

	// User Management

	secured.Post("/users", userHandler.RegisterUser)
	secured.Put("/users/:id", userHandler.UpdateUser)
	secured.Patch("/users/:id/status", userHandler.UpdateUserStatus)
	secured.Get("/users/:id", userHandler.GetUserByID)
	secured.Get("/users", userHandler.ListUsers)

	// Shift Management

	secured.Get("/shifts", shiftHandler.GetShiftByDate)
	secured.Get("/shifts/current", shiftHandler.GetCurrentShift)
	secured.Post("/shifts/open", shiftHandler.OpenShift)
	secured.Get("/shifts/:id", shiftHandler.GetShiftByID)
	secured.Patch("/shifts/:id/close", shiftHandler.CloseShift)


	// Ticket Management

	secured.Get("/tickets/:id", ticketHandler.GetTicket)
	secured.Get("/tickets/shifts/:shiftId", ticketHandler.GetTicketByShift)
	secured.Post("/tickets/sell", ticketHandler.CreateTicket)
	secured.Patch("/tickets/cancel/:code", ticketHandler.CancelTicket)
	secured.Patch("/tickets/use/:code", ticketHandler.UseTicket)

	// Ticket type Management
	secured.Get("/ticket-types", ticketTypeHandler.ListTicketType)
	secured.Post("/ticket-types", ticketTypeHandler.CreatedType)
	secured.Get("/ticket-types/:id", ticketTypeHandler.GetTicketType)
	secured.Put("/ticket-types/:id", ticketTypeHandler.UpdateTicketType)
	secured.Patch("/ticket-types/:id", ticketTypeHandler.UpdateTicketTypeStatus)

	// ==========================================
	// START APPLICATION SERVER
	// ==========================================
	log.Printf("Application starting in %s mode on port %s", cfg.AppEnv, cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
