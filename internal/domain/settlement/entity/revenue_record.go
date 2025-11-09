package entity

import "time"

type RevenueRecord struct {
	ID         string
	BookingID  string
	DeltaCents int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
