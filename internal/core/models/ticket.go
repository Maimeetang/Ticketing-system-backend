package domains

import "time"

type TicketStatus string

const (
	TicketActive    TicketStatus = "ACTIVE"
	TicketUsed      TicketStatus = "USED"
	TicketCancelled TicketStatus = "CANCELLED"
)

type Ticket struct {
	ID 		uint `gorm:"primaryKey"`
	ShiftID uint `gorm:"not null;index"`

	TicketTypeID 	uint	`gorm:"not null"`
	TicketTypeName 	string `gorm:"not null"`
	TicketCode 	 	string `gorm:"unique;not null"`

	Quantity uint `gorm:"not null"`

	UnitPrice  int64 `gorm:"not null"`
	TotalPrice int64 `gorm:"not null"`

	Status TicketStatus `gorm:"not null"`

	SoldAt 		time.Time  `gorm:"not null"`
	UsedAt 		*time.Time
	CancelledAt *time.Time 
}
