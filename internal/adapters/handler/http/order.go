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
	service port.OrderService
}

func NewOrderHandler(service port.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// POST /Orders/:id

// dto
type CreateOrderRequest struct {
	PaymentMethod domain.PaymentMethod  `json:"payment_method" validate:"required"`
	Tickets       []CreateTicketRequest `json:"tickets" validate:"required,dive"`
}

type CreateTicketRequest struct {
	TicketInfo []CreateTicketInfoRequest `json:"ticket_info" validate:"required,dive"`
}

type CreateTicketInfoRequest struct {
	TicketTypeID uint `json:"ticket_type_id" validate:"gt=0"`
	Quantity     int  `json:"quantity" validate:"gt=0"`
}

// CreateOrder
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	// id from cookie
	userIDCookie := c.Locals("user_id")
	if userIDCookie == nil {
		return apperror.NewUnauthorized("ไม่ได้รับอนุญาต: ไม่พบข้อมูลผู้ใช้งาน")
	}
	userID, ok := userIDCookie.(uint)
	if !ok {
		if userIDFloat, ok := userIDCookie.(float64); ok {
			userID = uint(userIDFloat)
		} else {
			return apperror.NewInternalServerError("ไม่สามารถแปลงประเภทข้อมูลผู้ใช้งานได้")
		}
	}

	// logic
	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.NewBadRequest("รูปแบบข้อมูลคำสั่งซื้อไม่ถูกต้อง")
	}

	if err := ValidateStruct(&req); err != nil {
		return err
	}

	order := &domain.Order{
		CashierID:     userID,
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

	createdOrder, err := h.service.CreateOrder(order)
	if err != nil {
		return err
	}

	res := newOrderResponse(createdOrder)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "สร้างรายการขายสำเร็จ",
		"data":    res,
	})
}

// GET /Orders/:id

// GetOrderByID
func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return apperror.NewBadRequest("รูปแบบรหัสคำสั่งซื้อไม่ถูกต้อง")
	}

	order, err := h.service.GetOrderByID(uint(id))

	if err != nil {
		return err
	}

	if order == nil {
		return apperror.NewNotFound("ไม่พบคำสั่งซื้อ")
	}

	res := newOrderResponse(order)

	return c.Status(fiber.StatusOK).JSON(res)
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
		return apperror.NewBadRequest("รูปแบบ query parameter ไม่ถูกต้อง")
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
			return apperror.NewBadRequest("รูปแบบวันที่ไม่ถูกต้อง (ต้องเป็น YYYY-MM-DD)")
		}
		filter.From = &parsedFrom
	}
	if req.To != "" {
		parsedTo, err := time.Parse("2006-01-02 15:04:05", req.To+" 23:59:59")
		if err != nil {
			return apperror.NewBadRequest("รูปแบบวันที่ไม่ถูกต้อง (ต้องเป็น YYYY-MM-DD)")
		}
		filter.To = &parsedTo
	}

	orders, totalCount, err := h.service.ListOrders(filter)
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
	res := NewPaginationResponse(records, totalCount, param, "/orders")

	return c.Status(fiber.StatusOK).JSON(res)
}

// PUT