package dto

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

type UpdateUserStatusRequest struct {
	IsActive *bool `json:"isActive" validate:"required"`
}

type CreateTicketBodyRequest struct {
	TicketTypeID uint `json:"ticketTypeId"`
	Quantity     uint `json:"quantity"`
}

type CancelledTicketBodyRequest struct {
	Remarks string `json:"remark"`
}

type TypeReq struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}