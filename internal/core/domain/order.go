package domain

import "time"

type Order struct {
	ID 			uint
	UserID 		uint
	TotalAmount float32
	CreatedAt 	time.Time
}