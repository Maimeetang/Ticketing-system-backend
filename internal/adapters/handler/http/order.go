package http

import (
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	orderService port.OrderService
}

func NewOrderHandler(orderService port.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

type CreateOrderRequest struct {
	PaymentMethod domain.PaymentMethod `json:"payment_method" binding:"required"`
	Tickets       []CreateTicketRequest `json:"tickets" binding:"required,dive"`
}

type CreateTicketRequest struct {
	TicketInfo []CreateTicketInfoRequest `json:"ticket_info" binding:"required,dive"`
}

type CreateTicketInfoRequest struct {
	TicketTypeID uint `json:"ticket_type_id" binding:"required,gt=0"`
	Quantity     int  `json:"quantity" binding:"required,gt=0"`
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	// id from cookie
	userIDCookie := c.Locals("user_id")
	shiftIDCookie := c.Locals("shift_id")
	if userIDCookie == nil || shiftIDCookie == nil {
		return apperror.NewUnauthorized("unauthorized: identity context missing")
	}
	userID, ok := userIDCookie.(uint)
	if !ok {
		if userIDFloat, ok := userIDCookie.(float64); ok {
			userID = uint(userIDFloat)
		} else {
			return apperror.NewInternalServerError("failed to resolve user identity type")
		}
	}

	shiftID, ok := shiftIDCookie.(uint)
	if !ok {
		if shiftIDFloat, ok := shiftIDCookie.(float64); ok {
			shiftID = uint(shiftIDFloat)
		} else {
			return apperror.NewInternalServerError("failed to resolve shift tracking identity type")
		}
	}

	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid order request body"})
	}

	order := &domain.Order{
		CashierID:     userID,
		ShiftID:       shiftID,
		PaymentMethod: req.PaymentMethod,
	}

	for _, tReq := range req.Tickets {
		ticket := domain.Ticket{}
		for _, infoReq := range tReq.TicketInfo {
			ticket.TicketInfos = append(ticket.TicketInfos, domain.TicketInfo{
				TicketTypeID: infoReq.TicketTypeID,
				Quantity:     infoReq.Quantity,
			})
		}
		order.Tickets = append(order.Tickets, ticket)
	}

	createdOrder, err := h.orderService.CreateOrder(order)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "order created successfully",
		"data":    createdOrder,
	})
}

