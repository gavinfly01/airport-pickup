package entity

import (
	"errors"
	"time"
)

// DriverOffer represents a driver's offer to serve airport pickup.
type DriverOffer struct {
	ID            string
	DriverID      string
	AirportCode   string
	VehicleType   string
	AvailableFrom time.Time
	AvailableTo   time.Time
	PricePerKm    float64
	Rating        float64
	Status        string // open, matched, cancelled
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// MarkMatched 将报价状态从 open 变为 matched，仅允许 open->matched
func (o *DriverOffer) MarkMatched() error {
	if o.Status != "open" {
		return errors.New("driver offer status must be 'open' to mark as 'matched'")
	}
	o.Status = "matched"
	return nil
}

// MarkCompleted 将报价状态从 matched 变为 completed，仅允许 matched->completed
func (o *DriverOffer) MarkCompleted() error {
	if o.Status != "matched" {
		return errors.New("driver offer status must be 'matched' to mark as 'completed'")
	}
	o.Status = "completed"
	return nil
}
