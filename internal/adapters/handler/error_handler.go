package handler

import (
	"errors"

	coreerr "ticketing-system/internal/core/errors"

	"github.com/gofiber/fiber/v2"
)

func handleError(c *fiber.Ctx, err error) error {

	var notFound coreerr.NotFound
	if errors.As(err, &notFound) {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	var conflict coreerr.Conflict
	if errors.As(err, &conflict) {
		return c.Status(409).JSON(fiber.Map{"error": err.Error()})
	}

	var bad coreerr.BadRequest
	if errors.As(err, &bad) {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(500).JSON(fiber.Map{
		"error": "internal server error",
	})
}