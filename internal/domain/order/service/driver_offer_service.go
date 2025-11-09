package service

import (
	"errors"
	"time"

	orderentity "github.com/gavin/airport-pickup/internal/domain/order/entity"
)

// DriverOfferService 负责创建司机报价领域对象
type DriverOfferService struct{}

// CreateDriverOfferCmd 封装创建司机报价的参数
type CreateDriverOfferCmd struct {
	DriverID      string
	AirportCode   string
	VehicleType   string
	AvailableFrom string
	AvailableTo   string
	PricePerKm    float64
	Rating        float64
}

// CreateDriverOffer 校验输入并创建司机报价领域对象（不生成ID，由上层或仓库负责）
func (s *DriverOfferService) CreateDriverOffer(cmd *CreateDriverOfferCmd) (*orderentity.DriverOffer, error) {
	if cmd.DriverID == "" {
		return nil, errors.New("driver_id required")
	}
	if cmd.AirportCode == "" {
		return nil, errors.New("airport_code required")
	}
	if cmd.VehicleType == "" {
		return nil, errors.New("vehicle_type required")
	}
	from, err := time.Parse(time.RFC3339, cmd.AvailableFrom)
	if err != nil {
		return nil, errors.New("invalid available_from")
	}
	to, err := time.Parse(time.RFC3339, cmd.AvailableTo)
	if err != nil {
		return nil, errors.New("invalid available_to")
	}
	if !to.After(from) {
		return nil, errors.New("available_to must be after available_from")
	}
	if cmd.PricePerKm <= 0 {
		return nil, errors.New("price_per_km required")
	}
	if cmd.Rating < 0 || cmd.Rating > 5 {
		return nil, errors.New("invalid rating")
	}
	return &orderentity.DriverOffer{
		ID:            "",
		DriverID:      cmd.DriverID,
		AirportCode:   cmd.AirportCode,
		VehicleType:   cmd.VehicleType,
		AvailableFrom: from,
		AvailableTo:   to,
		PricePerKm:    cmd.PricePerKm,
		Rating:        cmd.Rating,
		Status:        "open",
	}, nil
}
