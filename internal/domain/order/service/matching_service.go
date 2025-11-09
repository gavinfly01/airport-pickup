package service

import (
	"errors"
	"sort"
	"time"

	order "github.com/gavin/airport-pickup/internal/domain/order"
	orderentity "github.com/gavin/airport-pickup/internal/domain/order/entity"
	user "github.com/gavin/airport-pickup/internal/domain/user"
)

type MatchingService interface {
	// MatchFromCandidates matches using provided candidates (e.g., from in-memory order book) without hitting repository.
	MatchFromCandidates(req *orderentity.PickupRequest, candidates []*orderentity.DriverOffer) (*orderentity.DriverOffer, error)
	// CreateBooking 根据请求和报价生成 Booking 领域对象
	CreateBooking(req *orderentity.PickupRequest, offer *orderentity.DriverOffer, idGen func() string) *orderentity.Booking
}

type matchingService struct {
	orderRepo order.OrderRepository
	userRepo  user.DriverRepository
}

func NewMatchingService(orderRepo order.OrderRepository, userRepo user.DriverRepository) MatchingService {
	return &matchingService{orderRepo: orderRepo, userRepo: userRepo}
}

func (s *matchingService) MatchFromCandidates(req *orderentity.PickupRequest, candidates []*orderentity.DriverOffer) (*orderentity.DriverOffer, error) {
	if req == nil {
		return nil, errors.New("nil request")
	}
	// 先按业务规则过滤
	filtered := make([]*orderentity.DriverOffer, 0, len(candidates))
	for _, o := range candidates {
		if o == nil {
			continue
		}
		if o.AirportCode != req.AirportCode || o.VehicleType != req.VehicleType {
			continue
		}
		if !timeInRange(req.DesiredTime, o.AvailableFrom, o.AvailableTo) {
			continue
		}
		if o.PricePerKm > req.MaxPricePerKm {
			continue
		}
		filtered = append(filtered, o)
	}
	return s.rankAndPick(req, filtered)
}

func (s *matchingService) rankAndPick(req *orderentity.PickupRequest, candidates []*orderentity.DriverOffer) (*orderentity.DriverOffer, error) {
	if len(candidates) == 0 {
		return nil, errors.New("no offer matched")
	}
	// Load ratings when missing
	type ranked struct {
		offer  *orderentity.DriverOffer
		rating float64
	}
	rankedList := make([]ranked, 0, len(candidates))
	for _, c := range candidates {
		r := c.Rating
		rankedList = append(rankedList, ranked{offer: c, rating: r})
	}
	if len(rankedList) == 0 {
		return nil, errors.New("no offer with driver info")
	}

	// Sort based on preference
	if req.PreferHighRating {
		sort.Slice(rankedList, func(i, j int) bool {
			if rankedList[i].rating == rankedList[j].rating {
				return rankedList[i].offer.PricePerKm < rankedList[j].offer.PricePerKm
			}
			return rankedList[i].rating > rankedList[j].rating
		})
	} else {
		sort.Slice(rankedList, func(i, j int) bool {
			if rankedList[i].offer.PricePerKm == rankedList[j].offer.PricePerKm {
				return rankedList[i].rating > rankedList[j].rating
			}
			return rankedList[i].offer.PricePerKm < rankedList[j].offer.PricePerKm
		})
	}

	return rankedList[0].offer, nil
}

// CreateBooking 根据请求和报价生成 Booking 领域对象
func (s *matchingService) CreateBooking(req *orderentity.PickupRequest, offer *orderentity.DriverOffer, idGen func() string) *orderentity.Booking {
	margin := req.MaxPricePerKm - offer.PricePerKm
	if margin < 0 {
		margin = 0
	}
	return &orderentity.Booking{
		ID:                  idGen(),
		RequestID:           req.ID,
		OfferID:             offer.ID,
		PassengerID:         req.PassengerID,
		DriverID:            offer.DriverID,
		PricePerKm:          offer.PricePerKm,
		PlatformMarginPerKm: margin,
		Status:              "created",
	}
}

func timeInRange(t, from, to time.Time) bool {
	if to.Before(from) { // overnight window, normalize
		return !t.Before(from) || !t.After(to)
	}
	return !t.Before(from) && !t.After(to)
}
