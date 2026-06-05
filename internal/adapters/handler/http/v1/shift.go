package v1

import (
	"ticketing-system/internal/adapters/handler/http/utils"
	"ticketing-system/internal/adapters/handler/http/v1/dto"
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ShiftHandler struct {
	shiftService port.ShiftService
}

func NewShiftHandler(shiftService port.ShiftService) *ShiftHandler {
	return &ShiftHandler{shiftService: shiftService}
}

func (h *ShiftHandler) ClockIn(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	shift := domain.Shift{
		UserID:    userID,
		StartAt:   time.Now(),
		EndAt:     nil,
		Status:    domain.ShiftOpen,
	}

	_, err = h.shiftService.ClockIn(&shift)
	if err != nil {
		return err
	}

	res := dto.NewShiftResponse(&shift)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "เริ่มกะการทำงานสำเร็จ",
		"data":    res,
	})
}

func (h *ShiftHandler) ClockOut(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	err = h.shiftService.ClockOut(userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "สิ้นสุดกะการทำงานสำเร็จ",
	})
}

func (h *ShiftHandler) GetActiveShift(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	shift, err := h.shiftService.GetUserActiveShift(userID)
	if err != nil {
		return err
	}

	if shift == nil {
		return apperror.NewNotFound("ไม่พบกะการทำงานที่กำลังเปิดอยู่")
	}

	res := dto.NewShiftResponse(shift)

	return c.Status(fiber.StatusOK).JSON(res)
}
