package http

import (
	"ticketing-system/internal/adapters/handler/dto"
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

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	// Get user id from cookie
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

	// Get shift id from cookie
	shiftIDLocal := c.Locals("shift_id")
	if shiftIDLocal == nil {
		return apperror.NewForbidden("access denied: active shift tracking metadata missing")
	}
	shiftID, ok := shiftIDLocal.(uint)
	if !ok {
		if shiftIDFloat, ok := shiftIDLocal.(float64); ok {
			shiftID = uint(shiftIDFloat)
		} else {
			return apperror.NewInternalServerError("failed to resolve shift tracking identity type")
		}
	}

	// Check request body
	var req dto.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid order request body"})
	}

	// Map DTO slice parameters into Core Domain structures
	var domainItems []domain.OrderItem
	for _, item := range req.Items {
		domainItems = append(domainItems, domain.OrderItem{
			TicketTypeID: item.TicketTypeID,
			Quantity:     item.Quantity,
		})
	}

	// Invoke transaction execution through Core Service layer
	order, err := h.orderService.CreateOrder(userID, shiftID, domain.PaymentMethod(req.PaymentMethod), domainItems)
	if err != nil {
		return err // Dispatched directly to Centralized Error Middleware
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "order transaction placed successfully.",
		"uuid":    order.UUID,
		"total":   order.TotalAmount,
	})
}
