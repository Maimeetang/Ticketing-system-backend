package domain

import "time"

type Shift struct {
	ID        	 	uint        `gorm:"primaryKey"`
	UserID    	 	uint        `gorm:"not null;index:idx_user_status"`
	OpenAt    	 	time.Time	`gorm:"not null"`
	CloseAt   	 	*time.Time
	Status    	 	ShiftStatus `gorm:"type:varchar(20);not null;index:idx_user_status;default:'OPEN'"`
	TotalOrders  	uint		`gorm:"not null;default:0"`
	PaidOrders 		uint		`gorm:"not null;default:0"`
	CancelledOrders uint		`gorm:"not null;default:0"`
	TotalRevenue 	int64		`gorm:"not null;default:0"`
	User			User		`gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type ShiftFilter struct {
	UserID *uint
	Status *ShiftStatus
	StartDate *time.Time
}
