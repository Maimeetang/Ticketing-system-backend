package http

import (
	"ticketing-system/internal/adapters/handler/dto"
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/port"

	"github.com/gofiber/fiber/v2"
)

type ShiftHandler struct {
	shiftService port.ShiftService
}

func NewShiftHandler(shiftService port.ShiftService) *ShiftHandler {
	return &ShiftHandler{shiftService: shiftService}
}

func (h *ShiftHandler) ClockIn(c *fiber.Ctx) error {
	userIDLocal := c.Locals("user_id")
	if userIDLocal == nil {
		return apperror.NewUnauthorized("unauthorized: identity context missing")
	}

	userID, ok := userIDLocal.(uint)
	if !ok {
		if userIDFloat, ok := userIDLocal.(float64); ok {
			userID = uint(userIDFloat)
		} else {
			return apperror.NewInternalServerError("failed to resolve user identity type")
		}
	}

	shift, err := h.shiftService.ClockIn(userID)
	if err != nil {
		return err
	}

	res := dto.ShiftResponse{
		ID:      shift.ID,
		UserID:  shift.UserID,
		StartAt: shift.StartAt,
		Status:  string(shift.Status),
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "shift session opened successfully.",
		"data":    res,
	})
}

func (h *ShiftHandler) ClockOut(c *fiber.Ctx) error {
	userIDLocal := c.Locals("user_id")
	if userIDLocal == nil {
		return apperror.NewUnauthorized("unauthorized: identity context missing")
	}

	userID, ok := userIDLocal.(uint)
	if !ok {
		if userIDFloat, ok := userIDLocal.(float64); ok {
			userID = uint(userIDFloat)
		} else {
			return apperror.NewInternalServerError("failed to resolve user identity type")
		}
	}

	err := h.shiftService.ClockOut(userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "shift session closed successfully.",
	})
}

func (h *ShiftHandler) GetActiveShift(c *fiber.Ctx) error {
	userIDLocal := c.Locals("user_id")
	if userIDLocal == nil {
		return apperror.NewUnauthorized("unauthorized: identity context missing")
	}

	userID, ok := userIDLocal.(uint)
	if !ok {
		if userIDFloat, ok := userIDLocal.(float64); ok {
			userID = uint(userIDFloat)
		} else {
			return apperror.NewInternalServerError("failed to resolve user identity type")
		}
	}

	shift, err := h.shiftService.GetActiveShift(userID)
	if err != nil {
		return err
	}

	res := dto.ShiftResponse{
		ID:      shift.ID,
		UserID:  shift.UserID,
		StartAt: shift.StartAt,
		Status:  string(shift.Status),
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
