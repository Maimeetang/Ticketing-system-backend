package server

import (
	"ticketing-system/internal/adapters/handler/http"
	"ticketing-system/internal/adapters/handler/http/middleware"
)

// RegisterRoutes orchestrates and mounts all HTTP endpoint mappings
func (s *FiberServer) RegisterRoutes(
	authHandler *http.AuthHandler,
	userHandler *http.UserHandler,
	shiftHandler *http.ShiftHandler,
	orderHandler *http.OrderHandler,
	ticketTypeHandler *http.TicketTypeHandler,
	ticketHandler *http.TicketHandler,
) {
	// Root API Group
	// ----------------------------------------------------
	// Public Authentication Routes (Open Access Grid)
	// ----------------------------------------------------
	auth := s.App.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	// v1.Use(middleware.AuthRequired(s.Cfg))

	// ----------------------------------------------------
	// Protected User Administration (Guarded by role)
	// ----------------------------------------------------
	userRoutes :=  s.App.Group("/users")
	userRoutes.Post("/", userHandler.RegisterUser)
	userRoutes.Put("/:id", userHandler.UpdateUser)
	userRoutes.Patch("/:id/disable", userHandler.DisableUser)
	userRoutes.Patch("/:id/enable", userHandler.EnableUser)
	userRoutes.Get("/:id", userHandler.GetUserByID)
	userRoutes.Get("/", userHandler.ListUsers)

	 s.App.Use(middleware.AuthRequired(s.Cfg))

	// ----------------------------------------------------
	// Protected Shift Session Pipelines
	// ----------------------------------------------------
	shiftRoutes :=  s.App.Group("/shifts")
	shiftRoutes.Post("/open", shiftHandler.OpenShift)
	shiftRoutes.Put("/:id/close", shiftHandler.CloseShift)
	shiftRoutes.Get("/current", shiftHandler.GetCurrentShift)
	shiftRoutes.Get("/:id", shiftHandler.GetShiftByID)
	shiftRoutes.Get("/", shiftHandler.ListShifts)

	// ----------------------------------------------------
	// Protected Orders & Sales Module
	// ----------------------------------------------------
	orderRoutes :=  s.App.Group("/orders")
	orderRoutes.Post("/", orderHandler.CreateOrder)
	orderRoutes.Get("/:id", orderHandler.GetOrderByID)
	orderRoutes.Get("/", orderHandler.ListOrders)

	// ----------------------------------------------------
	// Protected Ticket Type
	// ----------------------------------------------------
	typeRoutes :=  s.App.Group("/ticket/types")
	typeRoutes.Post("/", ticketTypeHandler.CreatedType)
	typeRoutes.Put("/:id", ticketTypeHandler.UpdateType)
	typeRoutes.Get("/:id", ticketTypeHandler.GetTicketType)
	typeRoutes.Get("/", ticketTypeHandler.ListTicketType)

	// ----------------------------------------------------
	// Protected Ticket
	// ----------------------------------------------------
	ticketRoutes :=  s.App.Group("/ticket")
	ticketRoutes.Post("/use/:code", ticketHandler.UseTicket)
	ticketRoutes.Post("/cancell/:code", orderHandler.CancelOrder)
}
