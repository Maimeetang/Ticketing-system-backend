package model

type OrderItem struct {
	ID           uint
	orderID      uint
	TicketTypeID uint
	Quantity     uint
	UnitPrice    float32
	SubTotal     float32
}