package entity

import (
	"errors"
	"time"
)

// Booking represents a matched order (成交单).
type Booking struct {
	ID                  string
	RequestID           string
	OfferID             string
	PassengerID         string
	DriverID            string
	PricePerKm          float64
	PlatformMarginPerKm float64
	Status              string // created, completed, cancelled
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// MarkCompleted 将订单状态从 created 变为 completed，仅允许 created->completed
func (b *Booking) MarkCompleted() error {
	if b.Status != "created" {
		return errors.New("booking status must be 'created' to mark as 'completed'")
	}
	b.Status = "completed"
	return nil
}
