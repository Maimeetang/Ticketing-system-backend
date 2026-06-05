package dto

import "ticketing-system/internal/core/domain"

type CreateUserRequest struct {
	Username           string `json:"username" validate:"required,min=4"`
	Password           string `json:"password" validate:"required,min=6"`
	Role               string `json:"role" validate:"required,oneof=CASHIER SCANNER"`
	FirstName          string `json:"first_name" validate:"required"`
	LastName           string `json:"last_name" validate:"required"`
	PhoneNumber        string `json:"phone_number" validate:"required"`
	ReservePhoneNumber string `json:"reserve_phone_number"`
}

type UpdateUserRequest struct {
	Username           string `json:"username" validate:"required,min=4"`
	Role               string `json:"role" validate:"required,oneof=CASHIER SCANNER"`
	FirstName          string `json:"first_name" validate:"required"`
	LastName           string `json:"last_name" validate:"required"`
	PhoneNumber        string `json:"phone_number" validate:"required"`
	ReservePhoneNumber string `json:"reserve_phone_number"`
}

type CreateOrderRequest struct {
	PaymentMethod domain.PaymentMethod `json:"payment_method" validate:"required"`
	TicketTypeID  uint                 `json:"ticket_type_id" validate:"gt=0"`
	Quantity      uint                 `json:"quantity" validate:"gt=0"`
}

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

type TypeReq struct {
	Name  		string		`json:"name" validate:"required"`
	Price 		float64		`json:"price" validate:"gt=0"`
	Description string		`json:"description"`
	IsActive	bool		`json:"is_active" validate:"required"`
}