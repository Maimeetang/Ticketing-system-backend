package middleware

import (
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("access_token")
		if tokenString == "" {
			return apperror.NewUnauthorized("ไม่ได้รับอนุญาต: ไม่พบ access token")
		}

		// Parse token and verify cryptographic signature against configured secret
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperror.NewUnauthorized("วิธีการเข้ารหัสของ token ไม่ถูกต้อง")
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return apperror.NewUnauthorized("ไม่ได้รับอนุญาต: token ไม่ถูกต้องหรือหมดอายุแล้ว")
		}

		// Extract verified metadata/claims payload
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return apperror.NewUnauthorized("ไม่ได้รับอนุญาต: โครงสร้างข้อมูลใน token ไม่ถูกต้อง")
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("username", claims["username"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}
