package service

import (
	"errors"
	settlemententity "github.com/gavin/airport-pickup/internal/domain/settlement/entity"
)

// SettlementService defines settlement record building logic.
type SettlementService interface {
	CreateSettlementRecord(cmd *CreateSettlementRecordCmd) (*settlemententity.SettlementRecord, error)
	CreateRevenueRecord(cmd *CreateRevenueRecordCmd) (*settlemententity.RevenueRecord, error)
}

type settlementService struct{}

func NewSettlementService() SettlementService { return &settlementService{} }

type CreateSettlementRecordCmd struct {
	BookingID            string
	DriverID             string
	PassengerID          string
	AmountCents          int64
	PlatformRevenueCents int64
}

type CreateRevenueRecordCmd struct {
	BookingID  string
	DeltaCents int64
}

func (s *settlementService) CreateSettlementRecord(cmd *CreateSettlementRecordCmd) (*settlemententity.SettlementRecord, error) {
	if cmd.BookingID == "" || cmd.DriverID == "" || cmd.PassengerID == "" {
		return nil, errors.New("booking_id, driver_id, passenger_id required")
	}
	if cmd.AmountCents < 0 {
		return nil, errors.New("amount_cents must be >= 0")
	}
	if cmd.PlatformRevenueCents < 0 {
		return nil, errors.New("platform_revenue_cents must be >= 0")
	}
	return &settlemententity.SettlementRecord{
		ID:                   "",
		BookingID:            cmd.BookingID,
		DriverID:             cmd.DriverID,
		PassengerID:          cmd.PassengerID,
		AmountCents:          cmd.AmountCents,
		PlatformRevenueCents: cmd.PlatformRevenueCents,
	}, nil
}

func (s *settlementService) CreateRevenueRecord(cmd *CreateRevenueRecordCmd) (*settlemententity.RevenueRecord, error) {
	if cmd.BookingID == "" {
		return nil, errors.New("booking_id required")
	}
	if cmd.DeltaCents < 0 {
		return nil, errors.New("delta_cents must be >= 0")
	}
	return &settlemententity.RevenueRecord{
		ID:         "",
		BookingID:  cmd.BookingID,
		DeltaCents: cmd.DeltaCents,
	}, nil
}
