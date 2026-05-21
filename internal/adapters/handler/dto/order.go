package dto

import "time"

type CreateOrderItemRequest struct {
	TicketTypeID uint `json:"ticket_type_id" validate:"required"`
	Quantity     int  `json:"quantity" validate:"required,min=1"`
}

type CreateOrderRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required,oneof=cash epay"`
	Items []CreateOrderItemRequest `json:"items" validate:"required,div=1"`
}

type TicketResponse struct {
	TicketCode string `json:"ticket_code"`
	Status     string `json:"status"`
}

type OrderResponse struct {
	UUID        string           `json:"uuid"`
	TotalAmount float64          `json:"total_amount"`
	CreatedAt   time.Time        `json:"created_at"`
	Tickets     []TicketResponse `json:"tickets,omitempty"`
}
