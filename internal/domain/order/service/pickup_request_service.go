package service

import (
	"errors"
	"time"

	orderentity "github.com/gavin/airport-pickup/internal/domain/order/entity"
)

// PickupRequestService 负责创建接机请求领域对象
type PickupRequestService struct{}

// CreatePickupRequestCmd 封装创建接机请求的参数
type CreatePickupRequestCmd struct {
	PassengerID      string
	AirportCode      string
	VehicleType      string
	DesiredTime      string
	MaxPricePerKm    float64
	PreferHighRating bool
}

// CreatePickupRequest 校验输入并创建接机请求领域对象（不生成ID，由上层或仓库负责）
func (s *PickupRequestService) CreatePickupRequest(cmd *CreatePickupRequestCmd) (*orderentity.PickupRequest, error) {
	if cmd.PassengerID == "" {
		return nil, errors.New("passenger_id required")
	}
	if cmd.AirportCode == "" {
		return nil, errors.New("airport_code required")
	}
	if cmd.VehicleType == "" {
		return nil, errors.New("vehicle_type required")
	}
	t, err := time.Parse(time.RFC3339, cmd.DesiredTime)
	if err != nil {
		return nil, errors.New("invalid desired_time")
	}
	if cmd.MaxPricePerKm <= 0 {
		return nil, errors.New("max_price_per_km required")
	}
	return &orderentity.PickupRequest{
		ID:               "",
		PassengerID:      cmd.PassengerID,
		AirportCode:      cmd.AirportCode,
		VehicleType:      cmd.VehicleType,
		DesiredTime:      t,
		MaxPricePerKm:    cmd.MaxPricePerKm,
		PreferHighRating: cmd.PreferHighRating,
		Status:           "open",
	}, nil
}
