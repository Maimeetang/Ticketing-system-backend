package domain

import "time"

// TicketType represents the configuration of tickets (loaded read-only on POS)
type TicketType struct {
	ID    uint    `gorm:"primaryKey" json:"id"`
	Name  string  `gorm:"not null" json:"name"`
	Price float64 `gorm:"type:decimal(10,2);not null" json:"price"`
}

type Order struct {
	ID            uint        	`gorm:"primaryKey" json:"id"`
	UUID          string      	`gorm:"type:varchar(36);uniqueIndex;not null" json:"uuid"` // Globally unique ID for cloud sync
	ShiftID       uint        	`gorm:"not null;index" json:"shift_id"`
	TotalAmount   float64     	`gorm:"type:decimal(10,2);not null" json:"total_amount"`
	PaymentMethod PaymentMethod	`gorm:"type:varchar(20);not null;default:'CASH'" json:"payment_method"`
	SyncStatus    SyncStatus  	`gorm:"type:varchar(20);not null;default:'PENDING';index" json:"sync_status"`
	CreatedBy     uint        	`gorm:"not null;index" json:"created_by"`
	CreatedAt     time.Time   	`json:"created_at"`
	UpdatedAt     time.Time   	`json:"updated_at"`
	OrderItems    []OrderItem 	`gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"order_items,omitempty"`
}

type OrderItem struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	OrderID      uint    `gorm:"not null;index" json:"order_id"`
	TicketTypeID uint    `gorm:"not null;index" json:"ticket_type_id"`
	Quantity     int     `gorm:"not null" json:"quantity"`
	UnitPrice    float64 `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Subtotal     float64 `gorm:"type:decimal(10,2);not null" json:"subtotal"`
}

type Ticket struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	OrderItemID uint         `gorm:"not null;index" json:"order_item_id"`
	TicketCode  string       `gorm:"type:varchar(100);uniqueIndex;not null" json:"ticket_code"` // Secure random alphanumeric string
	Status      TicketStatus `gorm:"type:varchar(20);not null;default:'UNUSED'" json:"status"`
	SyncStatus  SyncStatus   `gorm:"type:varchar(20);not null;default:'PENDING';index" json:"sync_status"`
	UsedAt      *time.Time   `gorm:"default:null" json:"used_at"`
	CancelledAt *time.Time   `gorm:"default:null" json:"cancelled_at"`
	CreatedAt   time.Time    `json:"created_at"`
}
