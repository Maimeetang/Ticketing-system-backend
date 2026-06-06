package server

import (
	"ticketing-system/internal/adapters/handler/http/auth"
	"ticketing-system/internal/adapters/handler/http/middleware"
	v1 "ticketing-system/internal/adapters/handler/http/v1"
)

// RegisterRoutes orchestrates and mounts all HTTP endpoint mappings
func (s *FiberServer) RegisterRoutes(
	authHandler *auth.AuthHandler,
	userHandler *v1.UserHandler,
	shiftHandler *v1.ShiftHandler,
	orderHandler *v1.OrderHandler,
	ticketTypeHandler *v1.TicketTypeHandler,
	ticketHandler *v1.TicketHandler,
) {
	// Root API Group
	// ----------------------------------------------------
	// Public Authentication Routes (Open Access Grid)
	// ----------------------------------------------------
	auth := s.App.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	v1 := s.App.Group("/v1")
	// v1.Use(middleware.AuthRequired(s.Cfg))

	// ----------------------------------------------------
	// Protected User Administration (Guarded by role)
	// ----------------------------------------------------
	userRoutes := v1.Group("/users")
	userRoutes.Post("/", userHandler.RegisterUser)
	userRoutes.Put("/:id", userHandler.UpdateUser)
	userRoutes.Patch("/:id/disable", userHandler.DisableUser)
	userRoutes.Patch("/:id/enable", userHandler.EnableUser)
	userRoutes.Get("/:id", userHandler.GetUserByID)
	userRoutes.Get("/", userHandler.ListUsers)

	v1.Use(middleware.AuthRequired(s.Cfg))

	// ----------------------------------------------------
	// Protected Shift Session Pipelines
	// ----------------------------------------------------
	shiftRoutes := v1.Group("/shifts")
	shiftRoutes.Post("/open", shiftHandler.OpenShift)
	shiftRoutes.Put("/:id/close", shiftHandler.CloseShift)
	shiftRoutes.Get("/current", shiftHandler.GetCurrentShift)
	shiftRoutes.Get("/:id", shiftHandler.GetShiftByID)
	shiftRoutes.Get("/", shiftHandler.ListShifts)

	// ----------------------------------------------------
	// Protected Orders & Sales Module
	// ----------------------------------------------------
	orderRoutes := v1.Group("/orders")
	orderRoutes.Post("/", orderHandler.CreateOrder)
	orderRoutes.Get("/:id", orderHandler.GetOrderByID)
	orderRoutes.Get("/", orderHandler.ListOrders)

	// ----------------------------------------------------
	// Protected Ticket Type
	// ----------------------------------------------------
	typeRoutes := v1.Group("/ticket/types")
	typeRoutes.Post("/", ticketTypeHandler.CreatedType)
	typeRoutes.Put("/:id", ticketTypeHandler.UpdateType)
	typeRoutes.Get("/:id", ticketTypeHandler.GetTicketType)
	typeRoutes.Get("/", ticketTypeHandler.ListTicketType)

	// ----------------------------------------------------
	// Protected Ticket
	// ----------------------------------------------------
	ticketRoutes := v1.Group("/ticket")
	ticketRoutes.Post("/use/:code", ticketHandler.UseTicket)
	ticketRoutes.Post("/cancell/:code", orderHandler.CancelOrder)
}
