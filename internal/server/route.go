package server

import (
	"ticketing-system/internal/adapters/handler/http"
)

// RegisterRoutes orchestrates and mounts all HTTP endpoint mappings
func (s *FiberServer) RegisterRoutes(
	userHandler *http.UserHandler,
	authHandler *http.AuthHandler,
	shiftHandler *http.ShiftHandler,
) {
	// Root API Group
	api := s.App.Group("/api")

	// ----------------------------------------------------
	// Public Authentication Routes (Open Access Grid)
	// ----------------------------------------------------
	api.Post("/auth/login", authHandler.Login)
	api.Post("/auth/logout", authHandler.Logout)

	// ----------------------------------------------------
	// Protected User Administration (Guarded by JWTMiddleware)
	// ----------------------------------------------------
	userRoutes := api.Group("/users", http.JWTMiddleware(s.Cfg))
	userRoutes.Post("/", userHandler.Register)
	userRoutes.Put("/:id", userHandler.UpdateUser)
	userRoutes.Delete("/:id", userHandler.DeleteUser)
	userRoutes.Get("/:id", userHandler.GetUser)
	userRoutes.Get("/", userHandler.ListUsers)

	// ----------------------------------------------------
	// Protected Shift Session Pipelines (Guarded by JWTMiddleware)
	// ----------------------------------------------------
	shiftRoutes := api.Group("/shifts", http.JWTMiddleware(s.Cfg))
	shiftRoutes.Post("/clock-in", shiftHandler.ClockIn)
	shiftRoutes.Post("/clock-out", shiftHandler.ClockOut)
	shiftRoutes.Get("/active", shiftHandler.GetActiveShift)

	// ----------------------------------------------------
	// Protected Orders & Sales Module (Chained with Shift verification Gate)
	// ----------------------------------------------------
	// orderRoutes := api.Group("/orders", http.JWTMiddleware(s.Cfg), http.CheckActiveShift(shiftService))
	// orderRoutes.Post("/", orderHandler.CreateOrder)
}
