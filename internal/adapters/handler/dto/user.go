package dto

type UserResponse struct {
	ID                 uint   `json:"id"`
	Username           string `json:"username"`
	Role               string `json:"role"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	PhoneNumber        string `json:"phone_number"`
	ReservePhoneNumber string `json:"reserve_phone_number"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
	Role     string `json:"role"`
}