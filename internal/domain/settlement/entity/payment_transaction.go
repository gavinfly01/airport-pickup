package entity

import "time"

type PaymentTransaction struct {
	ID          string
	BookingID   string
	AmountCents int64
	Status      string // pending, success, failed
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
