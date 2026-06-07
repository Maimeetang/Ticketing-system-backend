package http

import (
	"strconv"
	"ticketing-system/internal/adapters/handler/http/dto"
	"ticketing-system/internal/adapters/handler/http/utils"
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ShiftHandler struct {
	service port.ShiftService
}

func NewShiftHandler(service port.ShiftService) *ShiftHandler {
	return &ShiftHandler{service: service}
}

// POST /shifts
func (h *ShiftHandler) OpenShift(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	shift, err := h.service.OpenShift(userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "เริ่มกะการทำงานสำเร็จ",
		"data":    dto.NewShiftResponse(shift),
	})
}

// POST /shifts/:id
func (h *ShiftHandler) CloseShift(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return apperror.NewBadRequest("id ผู้ใช้งานไม่ถูกต้อง")
	}

	err = h.service.CloseShift(uint(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "สิ้นสุดกะการทำงานสำเร็จ",
	})
}

// GET /shifts/current
func (h *ShiftHandler) GetCurrentShift(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	shift, err := h.service.GetCurrentShift(userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewShiftResponse(shift))
}

// GET /shifts/:id
func (h *ShiftHandler) GetShiftByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return apperror.NewBadRequest("id ผู้ใช้งานไม่ถูกต้อง")
	}

	shift, err := h.service.GetShiftByID(uint(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewShiftResponse(shift))
}

// GET /shifts
func (h *ShiftHandler) ListShifts(c *fiber.Ctx) error {
	var req dto.ListShiftRequest

	if err := c.QueryParser(&req); err != nil {
		return apperror.NewBadRequest("query parameters ไม่ถูกต้อง")
	}

	if req.StartDate == "" {
		return apperror.NewBadRequest("ต้องระบุ start_date ใน query parameters")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return apperror.NewBadRequest("start_date ต้องเป็นฟอร์แมต YYYY-MM-DD format")
	}

	filter := domain.ShiftFilter{
		UserID:    req.CashierID,
		Status:    req.Status,
		StartDate: &startDate,
	}

	shifts, err := h.service.ListShifts(filter)
	if err != nil {
		return err
	}

	res := make([]dto.ShiftResponse, 0, len(shifts))
	for _, shift := range shifts {
		res = append(res, dto.NewShiftResponse(&shift))
	}

	return c.JSON(res)
}