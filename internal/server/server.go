package server

import (
	"context"
	"log"
	"ticketing-system/internal/adapters/handler/http"
	"ticketing-system/internal/adapters/repository/orm"
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/config"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FiberServer struct {
	App *fiber.App
	Cfg *config.AuthConfig
}

func New() *FiberServer {
	// Initialize configuration values
	cfg := config.LoadAuthConfig()

	// Open internal SQLite file connection
	db, err := gorm.Open(sqlite.Open("./data/local_tickets.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("database initialization aborted (SQLite): %v", err)
	}

	// Execute AutoMigrate constraints to secure storage tables
	if err := db.AutoMigrate(
		&domain.User{}, 
		&domain.Shift{}, 
		&domain.Order{}, 
		&domain.Ticket{}, 
		&domain.TicketInfo{}, 
		&domain.TicketLog{},
		&domain.TicketType{},
	); err != nil {
		log.Fatalf("database schemas auto-migration failed: %v", err)
	}

	// Construct Hexagonal Dependency Injection pipeline matching ports contract
	userRepo := orm.NewGormUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := http.NewUserHandler(userService)

	authService := service.NewAuthService(userRepo, cfg)
	authHandler := http.NewAuthHandler(authService, cfg)

	shiftRepo := orm.NewGormShiftRepository(db)
	shiftService := service.NewShiftService(shiftRepo)
	shiftHandler := http.NewShiftHandler(shiftService)

	ticketTypeRepo := orm.NewGormTicketTypeRepository(db)
	ticketTypeService := service.NewTicketTypeService(ticketTypeRepo)
	ticketTypeHandler := http.NewTicketTypeHandler(ticketTypeService)

	orderRepo := orm.NewGormOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, ticketTypeRepo)
	orderHandler := http.NewOrderHandler(orderService, shiftService)

	// Build Fiber application engine with centralized error interceptor
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			statusCode := fiber.StatusInternalServerError
			message := "internal server processing fault occurred"

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

	// Create server instance reference
	server := &FiberServer{
		App: app,
		Cfg: cfg,
	}

	server.RegisterRoutes(userHandler, authHandler, shiftHandler, orderHandler, ticketTypeHandler)

	return server
}

func (s *FiberServer) Listen(addr string) error {
	return s.App.Listen(addr)
}

func (s *FiberServer) ShutdownWithContext(ctx context.Context) error {
	return s.App.ShutdownWithContext(ctx)
}