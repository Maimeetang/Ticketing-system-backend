package http

import (
	"ticketing-system/internal/adapters/handler/dto"
	"ticketing-system/internal/config"
	"ticketing-system/internal/core/port"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService port.AuthService
	cfg *config.AuthConfig
}

func NewAuthHandler(authService port.AuthService,  cfg *config.AuthConfig) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg: cfg,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		return err
	}

	// HttpOnly Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	})

	return c.Status(fiber.StatusOK).JSON(dto.LoginResponse{
		Message:  "Login success.",
		Username: req.Username,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout Success",
	})
}