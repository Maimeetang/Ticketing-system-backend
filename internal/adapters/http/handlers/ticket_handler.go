package handlers

import (
	"strconv"
	dto "ticketing-system/internal/adapters/http/dto"
	e "ticketing-system/internal/core/error"
	s "ticketing-system/internal/core/services"

	"github.com/gofiber/fiber/v2"
)

type TicketHandler struct {
	service s.TicketService
}

func NewTicketHandler(service s.TicketService) *TicketHandler{
	return &TicketHandler{service: service}
}

func (h *TicketHandler) CreateTicket(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var reqBody dto.CreateTicketBodyRequest
	ctx := c.UserContext()

	if err := c.BodyParser(&reqBody); err != nil {
		return e.NewBadRequest("Bad Request")
	}

	ticket, err := h.service.CreateTicket(
		ctx,userID,reqBody.TicketTypeID,reqBody.Quantity,
	)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "New ticket created",
		"data": dto.NewTicketResponse(ticket),
	})
}

func (h *TicketHandler) GetTicket(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return e.NewBadRequest("id must be a number")
	}

	ctx := c.UserContext()
	
	ticket, err := h.service.GetTicketByID(
		ctx, uint(id),
	)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.NewTicketResponse(ticket))
}

func (h *TicketHandler) UseTicket(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	tCode := c.Params("code")

	result, err := h.service.UseTicket(
		c.Context(), userID, tCode,
	)

	return c.Status(fiber.StatusCreated).JSON(
		dto.NewTicketUpdateResponse(result),
	)
}

func (h *TicketHandler) CancelTicket(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	tCode := c.Params("code")

	var reqBody dto.CancelledTicketBodyRequest

	if err := c.BodyParser(&reqBody); err != nil {
		return e.NewBadRequest("Bad Request")
	}

	result, err := h.service.CancelTicket(
		c.Context(), userID, tCode, reqBody.Remarks,
	)

	return c.Status(fiber.StatusCreated).JSON(
		dto.NewTicketUpdateResponse(result),
	)
}