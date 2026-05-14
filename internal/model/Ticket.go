package model

import "time"

type Ticket struct {
	ID          uint
	OrderItemID uint
	TicketCode  string
	Status      string
	UsedAt      time.Time
	CancelledAt time.Time
	CreatedAt   time.Time
}