package eventbus

import "time"

// Event is a domain event marker.
type Event interface {
	Name() string
}

// EventBus defines pub/sub interface for domain events.
type EventBus interface {
	Publish(evt Event)
	Subscribe(eventName string, handler func(Event))
}

// Common domain events
const (
	EventOrderMatched         = "OrderMatched"
	EventOrderCompleted       = "OrderCompleted"
	EventPaymentSucceeded     = "PaymentSucceeded"
	EventSettlementCreated    = "SettlementCreated"
	EventRevenueUpdated       = "RevenueUpdated"
	EventPickupRequestCreated = "PickupRequestCreated"
	EventDriverOfferCreated   = "DriverOfferCreated"
)

// OrderMatched payload
// Emitted when an order is matched between a pickup request and a driver offer.
type OrderMatched struct {
	BookingID     string
	RequestID     string
	DriverOfferID string
}

func (e OrderMatched) Name() string { return EventOrderMatched }

// OrderCompleted payload
// Emitted when a booking is completed.
type OrderCompleted struct {
	BookingID string
}

func (e OrderCompleted) Name() string { return EventOrderCompleted }

// PaymentSucceeded payload
// Emitted when payment succeeds for a booking.
type PaymentSucceeded struct {
	BookingID   string
	AmountCents int64
}

func (e PaymentSucceeded) Name() string { return EventPaymentSucceeded }

// SettlementCreated payload
// Emitted when settlement is created for a booking.
type SettlementCreated struct {
	BookingID string
}

func (e SettlementCreated) Name() string { return EventSettlementCreated }

// RevenueUpdated payload
// Emitted when platform revenue is updated.
type RevenueUpdated struct {
	BookingID  string
	DeltaCents int64
}

func (e RevenueUpdated) Name() string { return EventRevenueUpdated }

// PickupRequestCreated payload
type PickupRequestCreated struct {
	RequestID        string
	PassengerID      string
	AirportCode      string
	VehicleType      string
	MaxPricePerKm    float64
	PreferHighRating bool
	DesiredTime      time.Time
	Status           string // open, matched, cancelled
}

func (e PickupRequestCreated) Name() string { return EventPickupRequestCreated }

// DriverOfferCreated payload
type DriverOfferCreated struct {
	OfferID       string
	DriverID      string
	AirportCode   string
	VehicleType   string
	AvailableFrom time.Time
	AvailableTo   time.Time
	PricePerKm    float64
	Rating        float64
	Status        string // open, matched, cancelled
}

func (e DriverOfferCreated) Name() string { return EventDriverOfferCreated }
