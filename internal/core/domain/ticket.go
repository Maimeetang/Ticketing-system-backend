package domain

import "time"

type Ticket struct {
	ID          uint         `gorm:"primaryKey"`
	OrderID     uint         `gorm:"uniqueIndex"`
	TicketCode  string       `gorm:"type:varchar(100);uniqueIndex;not null"`
	Status      TicketStatus `gorm:"type:varchar(20);not null;default:'UNUSED'"`
	TotalPrice  float64      `gorm:"type:decimal(10,2);not null"`
	TicketInfos []TicketInfo `gorm:"foreignKey:TicketID;constraint:OnDelete:CASCADE;"`
	TicketLogs  []TicketLog  `gorm:"foreignKey:TicketID;constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}