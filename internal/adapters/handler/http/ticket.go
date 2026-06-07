package http

import (
	"ticketing-system/internal/adapters/handler/http/utils"
	"ticketing-system/internal/core/port"

	"github.com/gofiber/fiber/v2"
)

type TicketHandler struct {
	service port.TicketService
}

func NewTicketHandler(service port.TicketService) *TicketHandler{
	return &TicketHandler{service: service}
}

func (h *TicketHandler) UseTicket(c *fiber.Ctx) error {
	code := c.Params("code")

	userID, err := utils.GetUserID(c)

	if err != nil {
		return err
	}

	err = h.service.UseTicket(code,userID)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ใช้งานตั๋วแล้ว",
	})
}