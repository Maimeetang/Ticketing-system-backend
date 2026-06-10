package handlers

import (
	e "ticketing-system/internal/core/error"

	"github.com/gofiber/fiber/v2"
)

func getUserID(c *fiber.Ctx) (uint, error) {
	userIDLocal := c.Locals("user_id")
	if userIDLocal == nil {
		return 0, e.NewUnauthorized("cookies: user id not found")
	}

	if userID, ok := userIDLocal.(uint); ok {
		return userID, nil
	}

	if userIDFloat, ok := userIDLocal.(float64); ok {
		return uint(userIDFloat), nil
	}

	return 0, e.NewInternalServerError("cookies: invalid user id")
}