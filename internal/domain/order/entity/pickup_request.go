package entity

import (
	"errors"
	"time"
)

// PickupRequest represents a passenger's request for airport pickup.
type PickupRequest struct {
	ID               string
	PassengerID      string
	AirportCode      string
	VehicleType      string
	DesiredTime      time.Time
	MaxPricePerKm    float64
	PreferHighRating bool
	Status           string // open, matched, cancelled
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// MarkMatched 将请求状态从 open 变为 matched，仅允许 open->matched
func (r *PickupRequest) MarkMatched() error {
	if r.Status != "open" {
		return errors.New("pickup request status must be 'open' to mark as 'matched'")
	}
	r.Status = "matched"
	return nil
}

// MarkCompleted 将请求状态从 matched 变为 completed，仅允许 matched->completed
func (r *PickupRequest) MarkCompleted() error {
	if r.Status != "matched" {
		return errors.New("pickup request status must be 'matched' to mark as 'completed'")
	}
	r.Status = "completed"
	return nil
}
