package http

import "github.com/gavin/airport-pickup/internal/app/dto"

// OrderApp is the application service contract the HTTP layer depends on.
type OrderApp interface {
	CreatePassenger(name string) (string, error)
	CreateDriver(name string, rating float64) (string, error)
	CreatePickupRequest(in dto.CreatePickupRequestInput) (string, error)
	CreateDriverOffer(in dto.CreateDriverOfferInput) (string, error)
	ListBookings() ([]dto.BookingDTO, error)
	CompleteBooking(id string) error
}

// SettlementApp is reserved for future HTTP endpoints (e.g., manual payment trigger).
type SettlementApp interface {
	TriggerPayment(bookingID string) error
}

// Handler groups HTTP handlers and holds references to app services.
type Handler struct {
	orderApp      OrderApp
	settlementApp SettlementApp
}
