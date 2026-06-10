package handlers

import (
	dto "ticketing-system/internal/adapters/http/handlers/dto"
	e "ticketing-system/internal/core/error"
	s "ticketing-system/internal/core/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService s.AuthService
}

func NewAuthHandler(authService s.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx, CookieSecure bool, JWTSecret string) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return e.NewBadRequest("invalid request")
	}

	ctx := c.UserContext()

	token, err := h.authService.Login(ctx, req.Username, req.Password, JWTSecret)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   CookieSecure,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Login สำเร็จ",
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx, CookieSecure bool) error {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   CookieSecure,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout สำเร็จ",
	})
}