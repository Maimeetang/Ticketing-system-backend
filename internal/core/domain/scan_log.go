package domain

import "time"

type ScanLog struct {
	ID        uint
	UserID    uint
	TicketID  uint
	Action    string
	ScannedAt time.Time
}