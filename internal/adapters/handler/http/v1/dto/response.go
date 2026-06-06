package dto

import (
	"ticketing-system/internal/core/domain"
	"time"
)

// userResponse represents a user response body
type UserResponse struct {
	ID                 uint   `json:"id"`
	Username           string `json:"username"`
	Role               string `json:"role"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	PhoneNumber        string `json:"phone_number"`
	ReservePhoneNumber string `json:"reserve_phone_number"`
	IsActive		   bool	  `json:"is_active"`
}

// newUserResponse is a helper function to create a user body for handling user data
func NewUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID: 	  			user.ID,
		Username: 			user.Username,
		Role: 				string(user.Role),
		FirstName: 			user.FirstName,
		LastName: 			user.LastName,
		PhoneNumber: 		user.PhoneNumber,
		ReservePhoneNumber: user.ReservePhoneNumber,
		IsActive:			user.IsActive,
	}
}

// shiftResponse represents a shift response body
type ShiftResponse struct {
	ID        uint            `json:"id"`
	UserID    uint            `json:"user_id"`
	StartAt   time.Time       `json:"start_at"`
	EndAt     *time.Time      `json:"end_at"` 
	Status    string          `json:"status"`
	Orders    []OrderResponse `json:"orders"` 
}

// newShiftResponse is a helper function to create a response body for handling shift data
func NewShiftResponse(shift *domain.Shift) ShiftResponse {
	orders := make([]OrderResponse, 0)
	if len(shift.Orders) > 0 {
		for _, o := range shift.Orders {
			orders = append(orders, NewOrderResponse(&o))
		}
	}

	return ShiftResponse{
		ID:      shift.ID,
		UserID:  shift.UserID,
		StartAt: shift.StartAt,
		EndAt:   shift.EndAt,
		Status:  string(shift.Status),
		Orders:  orders,
	}
}

// orderResponse represents a order response body
type OrderResponse struct {
	ID            uint             `json:"id"`
	UserID        uint             `json:"user_id"`
	ShiftID       uint             `json:"shift_id"`
	TotalPrice    float64          `json:"total_price"`
	PaymentMethod string           `json:"payment_method"`
	Status        string           `json:"status"`
	Ticket        TicketResponse   `json:"ticket"` 
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

// newOrderResponse is a helper function to create a response body for handling order data
func NewOrderResponse(order *domain.Order) OrderResponse {
	return OrderResponse{
		ID:            order.ID,
		UserID:        order.UserID,
		ShiftID:       order.ShiftID,
		TotalPrice:    order.TotalPrice,
		PaymentMethod: string(order.PaymentMethod),
		Status:        string(order.Status),
		Ticket:        NewTicketResponse(&order.Ticket),
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
	}
}

// ticketResponse represents a ticket response body
type TicketResponse struct {
	ID           uint                `json:"id"`
	OrderID      uint                `json:"order_id"`
	TicketCode   string              `json:"ticket_code"`
	Status       string              `json:"status"`
	TicketType   string				 `json:"ticket_type"`
	Quantity     uint				 `json:"quantity"`
	PricePerUnit float64			 `json:"price_per_unit"`
	TicketLogs   []TicketLogResponse `json:"ticket_logs"` 
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// newTicketResponse is a helper function to create a response body for handling ticket data
func NewTicketResponse(ticket *domain.Ticket) TicketResponse {
	logs := make([]TicketLogResponse, 0)
	if len(ticket.TicketLogs) > 0 {
		for _, log := range ticket.TicketLogs {
			logs = append(logs, NewTicketLogResponse(&log))
		}
	}

	return TicketResponse{
		ID:         	ticket.ID,
		OrderID:    	ticket.OrderID,
		TicketCode: 	ticket.TicketCode,
		Status:     	string(ticket.Status),
		TicketType:		ticket.TicketType,
		Quantity:		ticket.Quantity,
		PricePerUnit: 	ticket.PricePerUnit,
		TicketLogs: 	logs,
		CreatedAt:  	ticket.CreatedAt,
		UpdatedAt:  	ticket.UpdatedAt,
	}
}

// ticketLogResponse represents a ticket log response body
type TicketLogResponse struct {
	ID          uint      `json:"id"`
	UserID	    uint      `json:"user_id"`
	TicketID    uint      `json:"ticket_id"`
	FromStatus  string    `json:"from_status"` 
	ToStatus    string    `json:"to_status"`
	Remarks     string    `json:"remarks"` 
	CreatedAt   time.Time `json:"created_at"`
}

// newTicketLogResponse is a helper function to create a response body for handling ticket log data
func NewTicketLogResponse(log *domain.TicketLog) TicketLogResponse {
	var fromStatusStr string
	if log.FromStatus != nil {
		fromStatusStr = string(*log.FromStatus)
	}
	
	return TicketLogResponse{
		ID:          log.ID,
		UserID: 	 log.UserID,
		TicketID:    log.TicketID,
		FromStatus:  fromStatusStr,
		ToStatus:    string(log.ToStatus),
		Remarks:     log.Remarks,
		CreatedAt:   log.CreatedAt,
	}
}