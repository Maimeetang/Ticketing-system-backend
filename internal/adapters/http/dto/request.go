package dto

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateUserRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
}

type UpdateUserRequest struct {
	Username    string `json:"username"`
	Role        string `json:"role"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
}

type ListShiftRequest struct {
	CashierID *uint  `query:"cashierId"`
	Status    string `query:"status"`
	StartDate string `query:"startDate"`
}

type CreateTicketBodyRequest struct {
	TicketTypeID uint
	Quantity     uint
}

type CancelledTicketBodyRequest struct {
	Remarks string
}

// type CreateOrderRequest struct {
// 	PaymentMethod domains.PaymentMethod `json:"payment_method" validate:"required"`
// 	TicketTypeID  uint                 `json:"ticket_type_id" validate:"gt=0"`
// 	Quantity      uint                 `json:"quantity" validate:"gt=0"`
// }

// type ListOrdersRequest struct {
// 	Page           int    `query:"page"`
// 	Limit          int    `query:"limit"`
// 	Status         string `query:"status"`
// 	PaymentMethod  string `query:"payment_method"`
// 	CashierID      uint   `query:"cashier_id"`
// 	ShiftID        uint   `query:"shift_id"`
// 	IncludeTickets bool   `query:"include_tickets"`
// 	Sort           string `query:"sort"`
// 	From           string `query:"from"`
// 	To             string `query:"to"`
// }

type TypeReq struct {
	Name        string  `json:"name" validate:"required"`
	Price       float64 `json:"price" validate:"gt=0"`
	Description string  `json:"description"`
	IsActive    bool    `json:"is_active" validate:"required"`
}