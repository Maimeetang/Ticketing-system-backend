package http

import (
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/config"
	"ticketing-system/internal/core/port"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(cfg *config.AuthConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("access_token")
		if tokenString == "" {
			return apperror.NewUnauthorized("unauthorized: missing access token")
		}

		// Parse token and verify cryptographic signature against configured secret
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return apperror.NewUnauthorized("unauthorized: invalid or expired token")
		}

		// Extract verified metadata/claims payload
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return apperror.NewUnauthorized("unauthorized: invalid payload structure")
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("username", claims["username"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}

func CheckActiveShift(shiftService port.ShiftService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDLocal := c.Locals("user_id")
		if userIDLocal == nil {
			return apperror.NewUnauthorized("unauthorized: identity context missing")
		}

		userID, ok := userIDLocal.(uint)
		if !ok {
			if userIDFloat, ok := userIDLocal.(float64); ok {
				userID = uint(userIDFloat)
			} else {
				return apperror.NewInternalServerError("failed to resolve user identity type within shift gate")
			}
		}

		activeShift, err := shiftService.GetActiveShift(userID)
		if err != nil || activeShift == nil {
			return apperror.NewForbidden("access denied: you must clock-in to open a working shift session before performing this action")
		}

		c.Locals("shift_id", activeShift.ID)

		return c.Next()
	}
}