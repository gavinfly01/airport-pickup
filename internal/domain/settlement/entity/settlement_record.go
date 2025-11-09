package entity

import "time"

type SettlementRecord struct {
	ID                   string
	BookingID            string
	DriverID             string
	PassengerID          string
	AmountCents          int64
	PlatformRevenueCents int64
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
