package handlers

import (
	"strconv"
	dto "ticketing-system/internal/adapters/http/dto"
	e "ticketing-system/internal/core/error"
	s "ticketing-system/internal/core/services"

	"github.com/gofiber/fiber/v2"
)

type TicketTypeHandler struct {
	service s.TicketTypeService
}

func NewTicketTypeHandler(service s.TicketTypeService) *TicketTypeHandler {
	return &TicketTypeHandler{service: service}
}

func (h *TicketTypeHandler) CreatedType(c *fiber.Ctx) error {
	var req dto.TypeReq

	if err := c.BodyParser(&req); err != nil {
		return e.NewBadRequest("bad request")
	}

	ctx := c.UserContext()

	createdType, err := h.service.CreateTicketType(
		ctx,req.Name,int64(req.Price),req.Description,
	)

	if err != nil{
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "new ticket type created",
		"data":    dto.NewTicketTypeResponse(createdType),
	})
}

func (h *TicketTypeHandler) UpdateTicketType(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return e.NewBadRequest("invalid id param")
	}

	var req dto.TypeReq

	if err := c.BodyParser(&req); err != nil {
		return e.NewBadRequest("bad request")
	}

	ctx := c.UserContext()
	updatedType, err := h.service.UpdateTicketType(
		ctx,uint(id),req.Name,int64(req.Price),req.Description,
	)

	if err != nil{
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ticket type updated",
		"data":    dto.NewTicketTypeResponse(updatedType),
	})
}

func (h *TicketTypeHandler) UpdateTicketTypeStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return e.NewBadRequest("invalid id param")
	}

	var req dto.UpdateTypeStatusReq
	if err := validate.Struct(req); err != nil {
		return e.NewBadRequest(err.Error())
	}

	ctx := c.UserContext()

	if err = h.service.UpdateTicketTypeStatus(ctx,uint(id),*req.IsActive); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ticket type status updated",
	})
}

func (h *TicketTypeHandler) GetTicketType(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return e.NewBadRequest("invalid id param")
	}

	ctx := c.UserContext()
	ticketType, err := h.service.GetTicketType(ctx, uint(id))

	if err != nil{
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(dto.NewTicketTypeResponse(ticketType))
}

func (h *TicketTypeHandler) ListTicketType(c *fiber.Ctx) error {
	withDisableStr := c.Query("withDisable")

	var withDisableBool bool
	var err error

	if withDisableStr != "" {
		withDisableBool, err = strconv.ParseBool(withDisableStr)
		if err != nil {
			return e.NewBadRequest("invalid withDisable query param (must be true or false)")
		}
	} else {
		withDisableBool = false
	}

	ctx := c.UserContext()

	types, err := h.service.ListTicketType(ctx, withDisableBool)

	if err != nil{
		return err
	}

	typesList := make([]dto.TicketTypeResponse, len(types))
	for i := range types {
		typesList[i] = dto.NewTicketTypeResponse(&types[i])
	}

	return c.Status(fiber.StatusOK).JSON(typesList)
}