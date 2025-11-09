package dto

// CreatePickupRequestInput represents passenger's pickup request creation input.
type CreatePickupRequestInput struct {
	PassengerID      string  `json:"passenger_id"`
	AirportCode      string  `json:"airport_code"`
	VehicleType      string  `json:"vehicle_type"`
	DesiredTime      string  `json:"desired_time"` // RFC3339
	MaxPricePerKm    float64 `json:"max_price_per_km"`
	PreferHighRating bool    `json:"prefer_high_rating"`
}

// CreateDriverOfferInput represents driver offer creation input.
type CreateDriverOfferInput struct {
	DriverID      string  `json:"driver_id"`
	AirportCode   string  `json:"airport_code"`
	VehicleType   string  `json:"vehicle_type"`
	AvailableFrom string  `json:"available_from"` // RFC3339
	AvailableTo   string  `json:"available_to"`   // RFC3339
	PricePerKm    float64 `json:"price_per_km"`
}

// BookingDTO is a simplified read model for bookings.
type BookingDTO struct {
	ID                  string  `json:"id"`
	RequestID           string  `json:"request_id"`
	OfferID             string  `json:"offer_id"`
	PassengerID         string  `json:"passenger_id"`
	DriverID            string  `json:"driver_id"`
	PricePerKm          float64 `json:"price_per_km"`
	PlatformMarginPerKm float64 `json:"platform_margin_per_km"`
	Status              string  `json:"status"`
}
