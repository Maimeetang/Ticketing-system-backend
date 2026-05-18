package handler

import (
	"errors"

	"ticketing-system/internal/apperror"

	"github.com/gofiber/fiber/v2"
)

func handleError(c *fiber.Ctx, err error) error {
	var (
		notFound apperror.NotFound
		conflict apperror.Conflict
		badRequest apperror.BadRequest
	)

	switch {
	case errors.As(err, &notFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	case errors.As(err, &conflict):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	case errors.As(err, &badRequest):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
}
