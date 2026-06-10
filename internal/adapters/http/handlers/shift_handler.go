package handlers

import (
	"strconv"
	dto "ticketing-system/internal/adapters/http/dto"
	e "ticketing-system/internal/core/error"
	s "ticketing-system/internal/core/services"

	"github.com/gofiber/fiber/v2"
)

type ShiftHandler struct {
	service s.ShiftService
}

func NewShiftHandler(service s.ShiftService) *ShiftHandler {
	return &ShiftHandler{service: service}
}

func (h *ShiftHandler) OpenShift(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	ctx := c.UserContext()

	shift, err := h.service.OpenShift(ctx, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "new shift opened",
		"data":    dto.NewShiftResponse(shift),
	})
}

func (h *ShiftHandler) CloseShift(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return e.NewBadRequest("invalid id param")
	}

	ctx := c.UserContext()

	err = h.service.CloseShift(ctx, uint(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "shift closed",
	})
}

func (h *ShiftHandler) GetCurrentShift(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	ctx := c.UserContext()

	shift, err := h.service.GetCurrentShift(ctx, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewShiftResponse(shift))
}

func (h *ShiftHandler) GetShiftByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return e.NewBadRequest("invalid id param")
	}

	ctx := c.UserContext()

	shift, err := h.service.GetShiftByID(ctx, uint(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewShiftResponse(shift))
}

func (h *ShiftHandler) GetShiftByDate(c *fiber.Ctx) error {
	date := c.Params("date")

	ctx := c.UserContext()

	shift, err := h.service.GetShiftByDate(ctx,date)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewShiftResponse(shift))
}