package dto

import (
	m "ticketing-system/internal/core/models"
	s "ticketing-system/internal/core/services"
	"time"
)

type UserResponse struct {
	ID                 uint   `json:"id"`
	Username           string `json:"username"`
	Role               string `json:"role"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	PhoneNumber        string `json:"phoneNumber"`
	IsActive		   bool	  `json:"isActive"`
}

func NewUserResponse(user *m.User) UserResponse {
	return UserResponse{
		ID: 	  			user.ID,
		Username: 			user.Username,
		Role: 				string(user.Role),
		FirstName: 			user.FirstName,
		LastName: 			user.LastName,
		PhoneNumber: 		user.PhoneNumber,
		IsActive:			user.IsActive,
	}
}

type ShiftResponse struct {
	ID        		 uint		   `json:"id"`
	UserID    		 uint		   `json:"userId"`
	User		 	 *UserResponse `json:"user,omitempty"`	
	OpenAt   		 time.Time	   `json:"openAt"`
	CloseAt     	 *time.Time	   `json:"closeAt"` 
	Status    		 string		   `json:"status"`
	TotalTickets	 uint		   `json:"totalTickets"`
	CancelledTickets uint		   `json:"cancelledTickets"`
	TotalRevenue	 int64		   `json:"totalRevenue"`
}

func NewShiftResponse(shift *m.Shift) ShiftResponse {
	var userResp *UserResponse

	if shift.User.ID != 0 {
		resp := NewUserResponse(&shift.User)
		userResp = &resp
	}

	return ShiftResponse{
		ID:      	 	  shift.ID,
		UserID:  	 	  shift.UserID,
		User: 	  		  userResp,
		OpenAt: 	 	  shift.OpenAt,
		CloseAt:   		  shift.CloseAt,
		Status:  		  string(shift.Status),
		TotalTickets: 	  shift.TotalTickets,
		CancelledTickets: shift.CancelledTickets,
		TotalRevenue: 	  shift.TotalRevenue,
	}
}

type TicketResponse struct {
	ID 				uint 		`json:"id"`
	ShiftID 		uint 		`json:"shiftId"`
	TicketTypeID 	uint		`json:"ticketTypeID"`
	TicketCode 	 	string 		`json:"ticketCode"`
	Quantity 		uint 		`json:"quantity"`
	UnitPrice  		int64 		`json:"unitPrice"`
	TotalPrice 		int64 		`json:"totalPrice"`
	Status 			string 		`json:"status"`
	SoldAt 			time.Time  	`json:"soldAt"`
	UsedAt 			*time.Time 	`json:"usedAt"`
	CancelledAt 	*time.Time 	`json:"cancelledAt"`
}

func NewTicketResponse(ticket *m.Ticket) TicketResponse {
	return TicketResponse{
		ID:           ticket.ID,
		ShiftID:      ticket.ShiftID,
		TicketTypeID: ticket.TicketTypeID,
		TicketCode:   ticket.TicketCode,
		Quantity:	  ticket.Quantity,
		UnitPrice:	  ticket.UnitPrice,
		TotalPrice:   ticket.TotalPrice,
		Status: 	  string(ticket.Status),
		SoldAt:  	  ticket.SoldAt,
		UsedAt:  	  ticket.UsedAt,
		CancelledAt:  ticket.CancelledAt,
	}
}

type TicketUpdateResponse struct {
	Status string
	Message string
	Ticket TicketResponse
}

func NewTicketUpdateResponse(result *s.StatusUpdateResult) TicketUpdateResponse {
	var message string

	switch result.Status {
	case s.ScanNotFound:
		message = "No ticket found"
	case s.ScanCancelled:
		message = "ticket has already been cancelled"
	case s.ScanUsed:
		message = "ticket has already been used"
	case s.ScanSuccess:
		message = "Ticket cancellation successful"
	}

	return TicketUpdateResponse{
		Status: string(result.Status),
		Message: message,
		Ticket: NewTicketResponse(result.Ticket),
	}
}

type TicketLogResponse struct {
	ID          uint      `json:"id"`
	UserID	    uint      `json:"userId"`
	TicketID    uint      `json:"ticketId"`
	FromStatus  string    `json:"fromStatus"` 
	ToStatus    string    `json:"toStatus"`
	Remarks     string    `json:"remarks"` 
	CreatedAt   time.Time `json:"createdAt"`
}

type TicketTypeResponse struct {
	ID 			uint		`json:"id"`
	Name 		string		`json:"name"`
	Price 		int64		`json:"price"`
	Description string		`json:"description"`
	IsActive 	bool		`json:"isActive"`
	CreatedAt 	time.Time	`json:"createdAt"`
	UpdatedAt 	time.Time	`json:"updatedAt"`
}

func NewTicketTypeResponse(tType *m.TicketType) TicketTypeResponse {
	return TicketTypeResponse{
		ID: tType.ID,
		Name: tType.Name,
		Price: tType.Price,
		Description: tType.Description,
		IsActive: tType.IsActive,
		CreatedAt: tType.CreatedAt,
		UpdatedAt: tType.UpdatedAt,
	}
}

// // newTicketLogResponse is a helper function to create a response body for handling ticket log data
// func NewTicketLogResponse(log *domains.TicketLog) TicketLogResponse {
// 	var fromStatusStr string
// 	if log.FromStatus != nil {
// 		fromStatusStr = string(*log.FromStatus)
// 	}
	
// 	return TicketLogResponse{
// 		ID:          log.ID,
// 		UserID: 	 log.UserID,
// 		TicketID:    log.TicketID,
// 		FromStatus:  fromStatusStr,
// 		ToStatus:    string(log.ToStatus),
// 		Remarks:     log.Remarks,
// 		CreatedAt:   log.CreatedAt,
// 	}
// }