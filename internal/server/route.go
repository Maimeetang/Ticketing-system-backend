package server

import (
	"ticketing-system/internal/adapters/handler/http"
)

// RegisterRoutes orchestrates and mounts all HTTP endpoint mappings
func (s *FiberServer) RegisterRoutes(
	userHandler *http.UserHandler,
	authHandler *http.AuthHandler,
	shiftHandler *http.ShiftHandler,
	orderHandler *http.OrderHandler,
) {
	// Root API Group
	// ----------------------------------------------------
	// Public Authentication Routes (Open Access Grid)
	// ----------------------------------------------------
	auth := s.App.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	s.App.Use(http.AuthRequired(s.Cfg))

	// ----------------------------------------------------
	// Protected User Administration (Guarded by role)
	// ----------------------------------------------------
	userRoutes := s.App.Group("/users")
	userRoutes.Post("/", userHandler.Register)
	userRoutes.Put("/:id", userHandler.UpdateUser)
	userRoutes.Delete("/:id", userHandler.DeleteUser)
	userRoutes.Get("/:id", userHandler.GetUser)
	userRoutes.Get("/", userHandler.ListUsers)

	// ----------------------------------------------------
	// Protected Shift Session Pipelines
	// ----------------------------------------------------
	shiftRoutes := s.App.Group("/shifts")
	shiftRoutes.Post("/start", shiftHandler.ClockIn)
	shiftRoutes.Post("/end", shiftHandler.ClockOut)
	shiftRoutes.Get("/active", shiftHandler.GetActiveShift)

	// ----------------------------------------------------
	// Protected Orders & Sales Module
	// ----------------------------------------------------
	orderRoutes := s.App.Group("/orders")
	orderRoutes.Post("/", orderHandler.CreateOrder)
	orderRoutes.Get("/:id", orderHandler.GetOrderByID)
	orderRoutes.Get("/", orderHandler.ListOrders)
}
