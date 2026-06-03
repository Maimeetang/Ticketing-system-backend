package domain

import "time"

type Shift struct {
	ID        uint        `gorm:"primaryKey"`
	UserID    uint        `gorm:"not null;index:idx_user_status;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	StartAt   time.Time
	EndAt     *time.Time
	Status    ShiftStatus `gorm:"type:varchar(20);not null;index:idx_user_status;default:'OPEN'"`
	Orders    []Order     `gorm:"foreignKey:ShiftID;"`
}
