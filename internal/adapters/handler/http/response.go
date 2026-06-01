package http

import (
	"ticketing-system/internal/core/domain"
	"time"
)

// userResponse represents a user response body
type userResponse struct {
	ID                 uint   `json:"id"`
	Username           string `json:"username"`
	Role               string `json:"role"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	PhoneNumber        string `json:"phone_number"`
	ReservePhoneNumber string `json:"reserve_phone_number"`
}

// newUserResponse is a helper function to create a user body for handling user data
func newUserResponse(user *domain.User) userResponse {
	return userResponse{
		ID: 	  			user.ID,
		Username: 			user.Username,
		Role: 				string(user.Role),
		FirstName: 			user.FirstName,
		LastName: 			user.LastName,
		PhoneNumber: 		user.PhoneNumber,
		ReservePhoneNumber: user.ReservePhoneNumber,
	}
}

// shiftResponse represents a shift response body
type shiftResponse struct {
	ID        uint            `json:"id"`
	UserID    uint            `json:"user_id"`
	StartAt   time.Time       `json:"start_at"`
	EndAt     *time.Time      `json:"end_at"` 
	Status    string          `json:"status"`
	Orders    []orderResponse `json:"orders"` 
}

// newShiftResponse is a helper function to create a response body for handling shift data
func newShiftResponse(shift *domain.Shift) shiftResponse {
	orders := make([]orderResponse, 0)
	if len(shift.Orders) > 0 {
		for _, o := range shift.Orders {
			orders = append(orders, newOrderResponse(&o))
		}
	}

	return shiftResponse{
		ID:      shift.ID,
		UserID:  shift.UserID,
		StartAt: shift.StartAt,
		EndAt:   shift.EndAt,
		Status:  string(shift.Status),
		Orders:  orders,
	}
}

// orderResponse represents a order response body
type orderResponse struct {
	ID            uint             `json:"id"`
	CashierID     uint             `json:"cashier_id"`
	ShiftID       uint             `json:"shift_id"`
	TotalPrice    float64          `json:"total_price"`
	PaymentMethod string           `json:"payment_method"`
	Status        string           `json:"status"`
	Ticket        ticketResponse   `json:"ticket"` 
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

// newOrderResponse is a helper function to create a response body for handling order data
func newOrderResponse(order *domain.Order) orderResponse {
	return orderResponse{
		ID:            order.ID,
		CashierID:     order.CashierID,
		ShiftID:       order.ShiftID,
		TotalPrice:    order.TotalPrice,
		PaymentMethod: string(order.PaymentMethod),
		Status:        string(order.Status),
		Ticket:        newTicketResponse(&order.Ticket),
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
	}
}

// ticketResponse represents a ticket response body
type ticketResponse struct {
	ID          uint                 `json:"id"`
	OrderID     uint                 `json:"order_id"`
	TicketCode  string               `json:"ticket_code"`
	Status      string               `json:"status"`
	TotalPrice  float64              `json:"total_price"`
	TicketInfos []ticketInfoResponse `json:"ticket_info"` 
	TicketLogs  []ticketLogResponse  `json:"ticket_logs"` 
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// newTicketResponse is a helper function to create a response body for handling ticket data
func newTicketResponse(ticket *domain.Ticket) ticketResponse {
	infos := make([]ticketInfoResponse, 0)
	if len(ticket.TicketInfos) > 0 {
		for _, info := range ticket.TicketInfos {
			infos = append(infos, newTicketInfoResponse(&info))
		}
	}

	logs := make([]ticketLogResponse, 0)
	if len(ticket.TicketLogs) > 0 {
		for _, log := range ticket.TicketLogs {
			logs = append(logs, newTicketLogResponse(&log))
		}
	}

	return ticketResponse{
		ID:          ticket.ID,
		OrderID:     ticket.OrderID,
		TicketCode:  ticket.TicketCode,
		Status:      string(ticket.Status),
		TotalPrice:  ticket.TotalPrice,
		TicketInfos: infos,
		TicketLogs:  logs,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}
}

// ticketInfoResponse represents a ticket info response body
type ticketInfoResponse struct {
	ID           uint      `json:"id"`
	TicketID     uint      `json:"ticket_id"`
	TicketTypeID uint      `json:"ticket_type_id"`
	Quantity     int       `json:"quantity"`
	PricePerUnit float64   `json:"price_per_unit"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// newTicketInfoResponse is a helper function to create a response body for handling ticket info data
func newTicketInfoResponse(info *domain.TicketInfo) ticketInfoResponse {
	return ticketInfoResponse{
		ID:           info.ID,
		TicketID:     info.TicketID,
		TicketTypeID: info.TicketTypeID,
		Quantity:     info.Quantity,
		PricePerUnit: info.PricePerUnit,
		CreatedAt:    info.CreatedAt,
		UpdatedAt:    info.UpdatedAt,
	}
}

// ticketLogResponse represents a ticket log response body
type ticketLogResponse struct {
	ID          uint      `json:"id"`
	TicketID    uint      `json:"ticket_id"`
	FromStatus  string    `json:"from_status"` 
	ToStatus    string    `json:"to_status"`
	TriggeredBy uint      `json:"triggered_by"`
	Remarks     string    `json:"remarks"` 
	CreatedAt   time.Time `json:"created_at"`
}

// newTicketLogResponse is a helper function to create a response body for handling ticket log data
func newTicketLogResponse(log *domain.TicketLog) ticketLogResponse {
	var fromStatusStr string
	if log.FromStatus != nil {
		fromStatusStr = string(*log.FromStatus)
	}
	
	return ticketLogResponse{
		ID:          log.ID,
		TicketID:    log.TicketID,
		FromStatus:  fromStatusStr,
		ToStatus:    string(log.ToStatus),
		TriggeredBy: log.TriggeredBy,
		Remarks:     log.Remarks,
		CreatedAt:   log.CreatedAt,
	}
}