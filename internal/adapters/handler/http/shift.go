package http

import (
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

func getUserID(c *fiber.Ctx) (uint, error) {
	userIDLocal := c.Locals("user_id")
	if userIDLocal == nil {
		return 0, apperror.NewUnauthorized("ไม่ได้รับอนุญาต: ไม่พบข้อมูลผู้ใช้งาน")
	}

	if userID, ok := userIDLocal.(uint); ok {
		return userID, nil
	}

	if userIDFloat, ok := userIDLocal.(float64); ok {
		return uint(userIDFloat), nil
	}

	return 0, apperror.NewInternalServerError("ไม่สามารถแปลงประเภทข้อมูลผู้ใช้งานได้")
}

func (h *ShiftHandler) ClockIn(c *fiber.Ctx) error {
	userID, err := getUserID(c)
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

	res := newShiftResponse(&shift)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "เริ่มกะการทำงานสำเร็จ",
		"data":    res,
	})
}

func (h *ShiftHandler) ClockOut(c *fiber.Ctx) error {
	userID, err := getUserID(c)
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
	userID, err := getUserID(c)
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

	res := newShiftResponse(shift)

	return c.Status(fiber.StatusOK).JSON(res)
}
