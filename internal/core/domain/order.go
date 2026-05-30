package domain

import "time"

type Order struct {
	ID            uint        	`gorm:"primaryKey"`
	CashierID     uint        	`gorm:"not null;index"`
	ShiftID       uint        	 `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	TotalPrice    float64     	`gorm:"type:decimal(10,2);not null"`
	PaymentMethod PaymentMethod	`gorm:"type:varchar(20);not null;default:'CASH'"`
	Status 		  OrderStatus	`gorm:"type:varchar(20);not null;default:'PAID'"`
	Tickets       []Ticket   	`gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type OrderFilter struct {
	Status         *OrderStatus
	PaymentMethod  *PaymentMethod
	CashierID      *uint
	ShiftID        *uint
	From      	   *time.Time
	To	           *time.Time
	Page           int
	Limit          int
	IncludeTickets bool
	Sort           string
}