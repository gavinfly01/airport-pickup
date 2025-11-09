package mysqlrepo

import (
	"time"

	"gorm.io/gorm"
)

// GORM models

type Passenger struct {
	ID        string    `gorm:"primaryKey;size:64"`
	Name      string    `gorm:"size:200;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type Driver struct {
	ID        string    `gorm:"primaryKey;size:64"`
	Name      string    `gorm:"size:200;not null"`
	Rating    float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type PickupRequest struct {
	ID               string    `gorm:"primaryKey;size:64"`
	PassengerID      string    `gorm:"index:idx_pickup_passenger_status;size:64;not null"`
	AirportCode      string    `gorm:"size:10;not null"`
	VehicleType      string    `gorm:"size:50;not null"`
	DesiredTime      time.Time `gorm:"not null"`
	MaxPricePerKm    float64   `gorm:"not null"`
	PreferHighRating bool      `gorm:"not null"`
	Status           string    `gorm:"size:20;index:idx_pickup_passenger_status;not null"`
	CreatedAt        time.Time `gorm:"not null"`
	UpdatedAt        time.Time `gorm:"not null"`
}

type DriverOffer struct {
	ID            string    `gorm:"primaryKey;size:64"`
	DriverID      string    `gorm:"index:idx_offer_driver_status;size:64;not null"`
	AirportCode   string    `gorm:"size:10;not null"`
	VehicleType   string    `gorm:"size:50;not null"`
	AvailableFrom time.Time `gorm:"not null"`
	AvailableTo   time.Time `gorm:"not null"`
	PricePerKm    float64   `gorm:"not null"`
	Rating        float64   `gorm:"not null"`
	Status        string    `gorm:"size:20;index:idx_offer_driver_status;not null"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
}

type Booking struct {
	ID                  string    `gorm:"primaryKey;size:64"`
	RequestID           string    `gorm:"size:64;not null"`
	OfferID             string    `gorm:"size:64;not null"`
	PassengerID         string    `gorm:"size:64;not null"`
	DriverID            string    `gorm:"size:64;not null"`
	PricePerKm          float64   `gorm:"not null"`
	PlatformMarginPerKm float64   `gorm:"not null"`
	Status              string    `gorm:"size:20;not null"`
	CreatedAt           time.Time `gorm:"not null"`
	UpdatedAt           time.Time `gorm:"not null"`
}

type PaymentTransaction struct {
	ID          string    `gorm:"primaryKey;size:64"`
	BookingID   string    `gorm:"size:64;not null"`
	AmountCents int64     `gorm:"not null"`
	Status      string    `gorm:"size:20;not null"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

type SettlementRecord struct {
	ID                   string    `gorm:"primaryKey;size:64"`
	BookingID            string    `gorm:"size:64;not null"`
	DriverID             string    `gorm:"size:64;not null"`
	PassengerID          string    `gorm:"size:64;not null"`
	AmountCents          int64     `gorm:"not null"`
	PlatformRevenueCents int64     `gorm:"not null"`
	CreatedAt            time.Time `gorm:"not null"`
	UpdatedAt            time.Time `gorm:"not null"`
}

type RevenueRecord struct {
	ID         string    `gorm:"primaryKey;size:64"`
	BookingID  string    `gorm:"size:64;not null"`
	DeltaCents int64     `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
}

// AutoMigrate migrates all tables.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Passenger{}, &Driver{},
		&PickupRequest{}, &DriverOffer{}, &Booking{},
		&PaymentTransaction{}, &SettlementRecord{}, &RevenueRecord{},
	)
}
