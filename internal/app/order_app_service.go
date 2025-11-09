package app

import (
	"errors"

	"github.com/gavin/airport-pickup/internal/app/dto"
	evt "github.com/gavin/airport-pickup/internal/domain/eventbus"
	order "github.com/gavin/airport-pickup/internal/domain/order"
	orderservice "github.com/gavin/airport-pickup/internal/domain/order/service"
	user "github.com/gavin/airport-pickup/internal/domain/user"
	userservice "github.com/gavin/airport-pickup/internal/domain/user/service"
	"github.com/gavin/airport-pickup/pkg/util"
)

type OrderAppService struct {
	orderRepo  order.OrderRepository
	passRepo   user.PassengerRepository
	driverRepo user.DriverRepository
	matching   orderservice.MatchingService
	bus        evt.EventBus

	passengerService     *userservice.PassengerService
	driverService        *userservice.DriverService
	pickupRequestService *orderservice.PickupRequestService
	driverOfferService   *orderservice.DriverOfferService
}

func NewOrderAppService(orderRepo order.OrderRepository, passRepo user.PassengerRepository, driverRepo user.DriverRepository, matching orderservice.MatchingService, bus evt.EventBus) *OrderAppService {
	return &OrderAppService{
		orderRepo:            orderRepo,
		passRepo:             passRepo,
		driverRepo:           driverRepo,
		matching:             matching,
		bus:                  bus,
		passengerService:     &userservice.PassengerService{},
		driverService:        &userservice.DriverService{},
		pickupRequestService: &orderservice.PickupRequestService{},
		driverOfferService:   &orderservice.DriverOfferService{},
	}
}

func (a *OrderAppService) CreatePassenger(name string) (string, error) {
	cmd := &userservice.CreatePassengerCmd{Name: name}
	p, err := a.passengerService.CreatePassenger(cmd)
	if err != nil {
		return "", err
	}
	p.ID = util.NewID()
	return p.ID, a.passRepo.Save(p)
}

func (a *OrderAppService) CreateDriver(name string, rating float64) (string, error) {
	cmd := &userservice.CreateDriverCmd{Name: name, Rating: rating}
	d, err := a.driverService.CreateDriver(cmd)
	if err != nil {
		return "", err
	}
	d.ID = util.NewID()
	return d.ID, a.driverRepo.Save(d)
}

func (a *OrderAppService) CreatePickupRequest(in dto.CreatePickupRequestInput) (string, error) {
	// 使用仓库方法判断是否存在进行中的请求（status in: open, matched）
	if ok, err := a.orderRepo.HasOngoingPickupRequest(in.PassengerID); err == nil && ok {
		return "", errors.New("ongoing pickup request exists")
	}
	cmd := &orderservice.CreatePickupRequestCmd{
		PassengerID:      in.PassengerID,
		AirportCode:      in.AirportCode,
		VehicleType:      in.VehicleType,
		DesiredTime:      in.DesiredTime,
		MaxPricePerKm:    in.MaxPricePerKm,
		PreferHighRating: in.PreferHighRating,
	}
	req, err := a.pickupRequestService.CreatePickupRequest(cmd)
	if err != nil {
		return "", err
	}
	req.ID = util.NewID()
	if err := a.orderRepo.SavePickupRequest(req); err != nil {
		return "", err
	}
	// 发布领域事件：创建接机请求
	a.bus.Publish(evt.PickupRequestCreated{RequestID: req.ID, PassengerID: req.PassengerID, AirportCode: req.AirportCode,
		VehicleType: req.VehicleType, MaxPricePerKm: req.MaxPricePerKm, PreferHighRating: req.PreferHighRating,
		DesiredTime: req.DesiredTime, Status: req.Status})
	return req.ID, nil
}

func (a *OrderAppService) CreateDriverOffer(in dto.CreateDriverOfferInput) (string, error) {
	// 使用仓库方法判断是否存在进行中的报价（status in: open, matched）
	if ok, err := a.orderRepo.HasOngoingDriverOffer(in.DriverID); err == nil && ok {
		return "", errors.New("ongoing driver offer exists")
	}
	driver, err := a.driverRepo.GetByID(in.DriverID)
	if err != nil {
		return "", err
	}
	if driver == nil {
		return "", errors.New("driver not found")
	}
	cmd := &orderservice.CreateDriverOfferCmd{
		DriverID:      in.DriverID,
		AirportCode:   in.AirportCode,
		VehicleType:   in.VehicleType,
		AvailableFrom: in.AvailableFrom,
		AvailableTo:   in.AvailableTo,
		PricePerKm:    in.PricePerKm,
		Rating:        driver.Rating,
	}
	o, err := a.driverOfferService.CreateDriverOffer(cmd)
	if err != nil {
		return "", err
	}
	o.ID = util.NewID()
	if err := a.orderRepo.SaveDriverOffer(o); err != nil {
		return "", err
	}
	// 发布领域事件：创建司机报价
	a.bus.Publish(evt.DriverOfferCreated{OfferID: o.ID, DriverID: o.DriverID, AirportCode: o.AirportCode, VehicleType: o.VehicleType,
		AvailableFrom: o.AvailableFrom, AvailableTo: o.AvailableTo, PricePerKm: o.PricePerKm, Rating: o.Rating, Status: o.Status})
	return o.ID, nil
}

func (a *OrderAppService) ListBookings() ([]dto.BookingDTO, error) {
	list, err := a.orderRepo.ListBookings()
	if err != nil {
		return nil, err
	}
	res := make([]dto.BookingDTO, 0, len(list))
	for _, b := range list {
		res = append(res, dto.BookingDTO{ID: b.ID, RequestID: b.RequestID, OfferID: b.OfferID, PassengerID: b.PassengerID, DriverID: b.DriverID, PricePerKm: b.PricePerKm, PlatformMarginPerKm: b.PlatformMarginPerKm, Status: b.Status})
	}
	return res, nil
}

func (a *OrderAppService) CompleteBooking(id string) error {
	b, err := a.orderRepo.GetBookingByID(id)
	if err != nil || b == nil {
		return errors.New("booking not found")
	}
	if err := b.MarkCompleted(); err != nil {
		return errors.New("booking mark completed failed: " + err.Error())
	}
	if b.RequestID == "" {
		return errors.New("booking.RequestID is empty")
	}
	if b.OfferID == "" {
		return errors.New("booking.OfferID is empty")
	}

	req, err := a.orderRepo.GetPickupRequestByID(b.RequestID)
	if err != nil {
		return errors.New("get pickup request failed: " + err.Error())
	}
	if req == nil {
		return errors.New("pickup request not found")
	}
	if err := req.MarkCompleted(); err != nil {
		return errors.New("pickup request mark completed failed: " + err.Error())
	}

	ofr, err := a.orderRepo.GetDriverOfferByID(b.OfferID)
	if err != nil {
		return errors.New("get driver offer failed: " + err.Error())
	}
	if ofr == nil {
		return errors.New("driver offer not found")
	}
	if err := ofr.MarkCompleted(); err != nil {
		return errors.New("driver offer mark completed failed: " + err.Error())
	}

	if err := a.orderRepo.UpdateAllInTransaction(b, req, ofr); err != nil {
		return err
	}
	a.bus.Publish(evt.OrderCompleted{BookingID: id})
	return nil
}
