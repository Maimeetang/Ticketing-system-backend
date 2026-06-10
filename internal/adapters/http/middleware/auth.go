package middleware

import (
	e "ticketing-system/internal/core/error"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(JWTSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("access_token")
		if tokenString == "" {
			return e.NewUnauthorized("unauthorized: access token not found")
		}

		// Parse token and verify cryptographic signature against configured secret
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, e.NewUnauthorized("Parse token fails")
			}
			return []byte(JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return e.NewUnauthorized("unauthorized: invalid token")
		}

		// Extract verified metadata/claims payload
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return e.NewUnauthorized("unauthorized: invalid token structure")
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("username", claims["username"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}
