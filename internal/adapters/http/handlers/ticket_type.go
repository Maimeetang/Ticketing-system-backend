package handlers

import (
	// "strconv"
	// "ticketing-system/internal/adapters/handler/http/dto"
	// "ticketing-system/internal/adapters/handler/http/validation"
	// "ticketing-system/internal/apperror"
	// "ticketing-system/internal/core/domains"
	s "ticketing-system/internal/core/services"
	// "github.com/gofiber/fiber/v2"
)

type TicketTypeHandler struct {
	service s.TicketTypeService
}

func NewTicketTypeHandler(service s.TicketTypeService) *TicketTypeHandler {
	return &TicketTypeHandler{service: service}
}

// POST /tickets/types
// func (h *TicketTypeHandler) CreatedType(c *fiber.Ctx) error {
// 	var req dto.TypeReq

// 	if err := c.BodyParser(&req); err != nil {
// 		return apperror.NewBadRequest("รูปแบบข้อมูลประเภทตั๋วไม่ถูกต้อง")
// 	}

// 	if err := validation.Validate(&req); err != nil {
// 		return err
// 	}

// 	ticketType := domains.TicketType{
// 		Name: req.Name,
// 		Price: req.Price,
// 		Description: req.Description,
// 		IsActive: req.IsActive,
// 	}

// 	createdType, err := h.service.CreateTicketType(&ticketType)

// 	if err != nil{
// 		return err
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"message": "สร้างประเภทตั๋วสำเร็จ",
// 		"data":    createdType,
// 	})
// }

// // PUT /tickets/types/:id
// func (h *TicketTypeHandler) UpdateType(c *fiber.Ctx) error {
// 	idParam := c.Params("id")
// 	id, err := strconv.ParseUint(idParam, 10, 32)
// 	if err != nil {
// 		return apperror.NewBadRequest("id ประเภทตั๋วไม่ถูกต้อง")
// 	}

// 	var req dto.TypeReq

// 	if err := c.BodyParser(&req); err != nil {
// 		return apperror.NewBadRequest("รูปแบบข้อมูลประเภทตั๋วไม่ถูกต้อง")
// 	}

// 	if err := validation.Validate(&req); err != nil {
// 		return err
// 	}

// 	ticketType := domains.TicketType{
// 		ID: uint(id),
// 		Name: req.Name,
// 		Price: req.Price,
// 		Description: req.Description,
// 		IsActive: req.IsActive,
// 	}

// 	updatedType, err := h.service.UpdateTicketType(&ticketType)

// 	if err != nil{
// 		return err
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "อัปเดตประเภทตั๋วสำเร็จ",
// 		"data":    updatedType,
// 	})
// }

// // GET /tickets/types/:id
// func (h *TicketTypeHandler) GetTicketType(c *fiber.Ctx) error {
// 	idParam := c.Params("id")
// 	id, err := strconv.ParseUint(idParam, 10, 32)
// 	if err != nil {
// 		return apperror.NewBadRequest("id ประเภทตั๋วไม่ถูกต้อง")
// 	}

// 	ticketType, err := h.service.GetTicketType(uint(id))

// 	if err != nil{
// 		return err
// 	}

// 	if ticketType == nil {
// 		return apperror.NewNotFound("ไม่พบประเภทตั๋ว")
// 	}

// 	return c.Status(fiber.StatusOK).JSON(ticketType)
// }

// // GET /tickets/types
// func (h *TicketTypeHandler) ListTicketType(c *fiber.Ctx) error {
// 	types, err := h.service.ListTicketType()

// 	if err != nil{
// 		return err
// 	}

// 	return c.Status(fiber.StatusOK).JSON(types)
// }