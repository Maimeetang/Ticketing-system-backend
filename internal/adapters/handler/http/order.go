package http

import (
	"strconv"
	"strings"
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/core/domain"
	"ticketing-system/internal/core/port"
	"time"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	orderService port.OrderService
	shiftService port.ShiftService
}

func NewOrderHandler(orderService port.OrderService, shiftService port.ShiftService) *OrderHandler {
	return &OrderHandler{orderService, shiftService}
}

// POST /Orders/:id

// dto
type CreateOrderRequest struct {
	PaymentMethod domain.PaymentMethod  `json:"payment_method" binding:"required"`
	Tickets       []CreateTicketRequest `json:"tickets" binding:"required,dive"`
}

type CreateTicketRequest struct {
	TicketInfo []CreateTicketInfoRequest `json:"ticket_info" binding:"required,dive"`
}

type CreateTicketInfoRequest struct {
	TicketTypeID uint `json:"ticket_type_id" binding:"required,gt=0"`
	Quantity     int  `json:"quantity" binding:"required,gt=0"`
}

// CreateOrder
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	// id from cookie
	userIDCookie := c.Locals("user_id")
	if userIDCookie == nil {
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

	// logic
	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid order request body"})
	}

	if err := ValidateStruct(&req); err != nil {
		return err
	}

	shift, err := h.shiftService.GetUserActiveShift(userID)
	if err != nil {
		return apperror.NewForbidden("you must clock-in first.")
	}

	order := &domain.Order{
		CashierID:     userID,
		ShiftID:       shift.ID,
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

// GET /Orders/:id

// GetOrderByID
func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return apperror.NewBadRequest("invalid order id format")
	}

	order, err := h.orderService.GetOrderByID(uint(id))
	if err != nil {
		return err
	}
	if order == nil {
		return apperror.NewNotFound("order not found.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": order,
	})
}

// GET /Orders

// dto
type ListOrdersRequest struct {
	Page           int    `query:"page"`
	Limit          int    `query:"limit"`
	Status         string `query:"status"`
	PaymentMethod  string `query:"payment_method"`
	CashierID      uint   `query:"cashier_id"`
	ShiftID        uint   `query:"shift_id"`
	IncludeTickets bool   `query:"include_tickets"`
	Sort           string `query:"sort"`
	From           string `query:"from"`
	To             string `query:"to"`
}

// helper function
func (h *OrderHandler) setListOrdersDefaults(f *ListOrdersRequest) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 10
	}
	if f.Sort == "" {
		f.Sort = "DESC"
	}
}

// ListOrders
func (h *OrderHandler) ListOrders(c *fiber.Ctx) error {
	var req ListOrdersRequest
	if err := c.QueryParser(&req); err != nil {
		return apperror.NewBadRequest("invalid query parameter format")
	}

	h.setListOrdersDefaults(&req)

	filter := domain.OrderFilter{
		Page:           req.Page,
		Limit:          req.Limit,
		IncludeTickets: req.IncludeTickets,
		Sort:           req.Sort,
	}

	if req.Status != "" {
		status := domain.OrderStatus(strings.ToUpper(req.Status))
		filter.Status = &status
	}
	if req.PaymentMethod != "" {
		method := domain.PaymentMethod(strings.ToUpper(req.PaymentMethod))
		filter.PaymentMethod = &method
	}
	if req.CashierID > 0 {
		filter.CashierID = &req.CashierID
	}
	if req.ShiftID > 0 {
		filter.ShiftID = &req.ShiftID
	}

	if req.From != "" {
		parsedFrom, err := time.Parse("2006-01-02", req.From)
		if err != nil {
			return apperror.NewBadRequest("invalid 'from' date format: must be YYYY-MM-DD")
		}
		filter.From = &parsedFrom
	}
	if req.To != "" {
		parsedTo, err := time.Parse("2006-01-02 15:04:05", req.To+" 23:59:59")
		if err != nil {
			return apperror.NewBadRequest("invalid 'to' date format: must be YYYY-MM-DD")
		}
		filter.To = &parsedTo
	}

	orders, totalCount, err := h.orderService.ListOrders(filter)
	if err != nil {
		return err
	}

	records := make([]orderResponse, 0, len(orders))
	for _, o := range orders {
		records = append(records, newOrderResponse(&o))
	}

	param := paginationParam{
		Page:    filter.Page,
		PerPage: filter.Limit,
	}
	response := NewPaginationResponse(records, totalCount, param, "/api/orders")
	return c.Status(fiber.StatusOK).JSON(response)
}

// PUT