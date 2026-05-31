package http

import (
	"strconv"
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"

	"github.com/gofiber/fiber/v2"
)

type TicketTypeHandler struct {
	service port.TicketTypeService
}

func NewTicketTypeHandler(service port.TicketTypeService) *TicketTypeHandler {
	return &TicketTypeHandler{service: service}
}

// dto
type TypeReq struct {
	Name  		string		`json:"name" validate:"required"`
	Price 		float64		`json:"price" validate:"gt=0"`
	Description string		`json:"description"`
	IsActive	bool		`json:"is_active" validate:"required"`
}

// POST /tickets/types
func (h *TicketTypeHandler) CreatedType(c *fiber.Ctx) error {
	var req TypeReq

	if err := c.BodyParser(&req); err != nil {
		return apperror.NewBadRequest("รูปแบบข้อมูลประเภทตั๋วไม่ถูกต้อง")
	}

	if err := ValidateStruct(&req); err != nil {
		return err
	}

	ticketType := domain.TicketType{
		Name: req.Name,
		Price: req.Price,
		Description: req.Description,
		IsActive: req.IsActive,
	}

	createdType, err := h.service.CreateTicketType(&ticketType)

	if err != nil{
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "สร้างประเภทตั๋วสำเร็จ",
		"data":    createdType,
	})
}

// PUT /tickets/types/:id
func (h *TicketTypeHandler) UpdateType(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return apperror.NewBadRequest("id ประเภทตั๋วไม่ถูกต้อง")
	}

	var req TypeReq

	if err := c.BodyParser(&req); err != nil {
		return apperror.NewBadRequest("รูปแบบข้อมูลประเภทตั๋วไม่ถูกต้อง")
	}

	if err := ValidateStruct(&req); err != nil {
		return err
	}

	ticketType := domain.TicketType{
		ID: uint(id),
		Name: req.Name,
		Price: req.Price,
		Description: req.Description,
		IsActive: req.IsActive,
	}

	updatedType, err := h.service.UpdateTicketType(&ticketType)

	if err != nil{
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "อัปเดตประเภทตั๋วสำเร็จ",
		"data":    updatedType,
	})
}

// GET /tickets/types/:id
func (h *TicketTypeHandler) GetTicketType(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return apperror.NewBadRequest("id ประเภทตั๋วไม่ถูกต้อง")
	}

	ticketType, err := h.service.GetTicketType(uint(id))

	if err != nil{
		return err
	}

	if ticketType == nil {
		return apperror.NewNotFound("ไม่พบประเภทตั๋ว")
	}

	return c.Status(fiber.StatusOK).JSON(ticketType)
}

// GET /tickets/types
func (h *TicketTypeHandler) ListTicketType(c *fiber.Ctx) error {
	types, err := h.service.ListTicketType()

	if err != nil{
		return err
	}

	return c.Status(fiber.StatusOK).JSON(types)
}