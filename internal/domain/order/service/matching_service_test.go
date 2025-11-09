package service

import (
	"github.com/gavin/airport-pickup/internal/domain/order/entity"
	"testing"
	"time"
)

type dummyIDGen struct{}

func (d dummyIDGen) Gen() string { return "id123" }

func TestMatchingService_MatchFromCandidates(t *testing.T) {
	svc := &matchingService{}
	req := &entity.PickupRequest{
		AirportCode:      "PVG",
		VehicleType:      "sedan",
		DesiredTime:      time.Date(2025, 11, 8, 10, 0, 0, 0, time.UTC),
		MaxPricePerKm:    10,
		PreferHighRating: true,
	}
	candidates := []*entity.DriverOffer{
		{
			ID: "1", AirportCode: "PVG", VehicleType: "sedan", AvailableFrom: time.Date(2025, 11, 8, 9, 0, 0, 0, time.UTC), AvailableTo: time.Date(2025, 11, 8, 12, 0, 0, 0, time.UTC), PricePerKm: 8, Rating: 4.9,
		},
		{
			ID: "2", AirportCode: "PVG", VehicleType: "sedan", AvailableFrom: time.Date(2025, 11, 8, 8, 0, 0, 0, time.UTC), AvailableTo: time.Date(2025, 11, 8, 11, 0, 0, 0, time.UTC), PricePerKm: 9, Rating: 4.7,
		},
		{
			ID: "3", AirportCode: "SHA", VehicleType: "sedan", AvailableFrom: time.Date(2025, 11, 8, 9, 0, 0, 0, time.UTC), AvailableTo: time.Date(2025, 11, 8, 12, 0, 0, 0, time.UTC), PricePerKm: 7, Rating: 4.8,
		},
		{
			ID: "4", AirportCode: "PVG", VehicleType: "van", AvailableFrom: time.Date(2025, 11, 8, 9, 0, 0, 0, time.UTC), AvailableTo: time.Date(2025, 11, 8, 12, 0, 0, 0, time.UTC), PricePerKm: 8, Rating: 4.6,
		},
	}
	best, err := svc.MatchFromCandidates(req, candidates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if best == nil || best.ID != "1" {
		t.Errorf("expected best offer ID '1', got %v", best)
	}

	// 测试无匹配
	req2 := &entity.PickupRequest{AirportCode: "PVG", VehicleType: "suv", DesiredTime: time.Date(2025, 11, 8, 10, 0, 0, 0, time.UTC), MaxPricePerKm: 10}
	best2, err2 := svc.MatchFromCandidates(req2, candidates)
	if err2 == nil {
		t.Errorf("expected error for no match, got nil")
	}
	if best2 != nil {
		t.Errorf("expected nil offer, got %v", best2)
	}
}

func TestMatchingService_CreateBooking(t *testing.T) {
	svc := &matchingService{}
	req := &entity.PickupRequest{ID: "req1", PassengerID: "p1", MaxPricePerKm: 10}
	offer := &entity.DriverOffer{ID: "off1", DriverID: "d1", PricePerKm: 8}
	idGen := func() string { return "bk1" }
	bk := svc.CreateBooking(req, offer, idGen)
	if bk == nil {
		t.Fatalf("expected booking, got nil")
	}
	if bk.ID != "bk1" || bk.RequestID != "req1" || bk.OfferID != "off1" || bk.PassengerID != "p1" || bk.DriverID != "d1" {
		t.Errorf("booking fields not set correctly: %+v", bk)
	}
	if bk.PlatformMarginPerKm != 2 {
		t.Errorf("expected margin 2, got %v", bk.PlatformMarginPerKm)
	}
	if bk.Status != "created" {
		t.Errorf("expected status 'created', got %v", bk.Status)
	}

	// 测试 margin < 0
	offer2 := &entity.DriverOffer{ID: "off2", DriverID: "d2", PricePerKm: 12}
	bk2 := svc.CreateBooking(req, offer2, idGen)
	if bk2.PlatformMarginPerKm != 0 {
		t.Errorf("expected margin 0, got %v", bk2.PlatformMarginPerKm)
	}
}
