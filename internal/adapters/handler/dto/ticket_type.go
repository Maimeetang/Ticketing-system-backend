package dto

// TicketTypeRequest acts as the boundary schema for creating/updating configurations
type TicketTypeRequest struct {
	Name  string  `json:"name" validate:"required"`
	Price float64 `json:"price" validate:"required,min=0"`
}

// TicketTypeResponse outputs sanitized configuration attributes back to admin dashboards
type TicketTypeResponse struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
