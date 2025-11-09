package order

import (
	orderentity "github.com/gavin/airport-pickup/internal/domain/order/entity"
)

type OrderRepository interface {
	// pickup requests
	SavePickupRequest(r *orderentity.PickupRequest) error
	GetPickupRequestByID(id string) (*orderentity.PickupRequest, error)
	ListPickupRequests() ([]*orderentity.PickupRequest, error)
	UpdatePickupRequest(r *orderentity.PickupRequest) error
	// 是否存在进行中的接机请求（status in: open, matched）
	HasOngoingPickupRequest(passengerID string) (bool, error)

	// driver offers
	SaveDriverOffer(o *orderentity.DriverOffer) error
	GetDriverOfferByID(id string) (*orderentity.DriverOffer, error)
	ListDriverOffers() ([]*orderentity.DriverOffer, error)
	UpdateDriverOffer(o *orderentity.DriverOffer) error
	// 是否存在进行中的司机报价（status in: open, matched）
	HasOngoingDriverOffer(driverID string) (bool, error)

	// bookings
	SaveBooking(b *orderentity.Booking) error
	GetBookingByID(id string) (*orderentity.Booking, error)
	ListBookings() ([]*orderentity.Booking, error)
	UpdateBooking(b *orderentity.Booking) error
	// 新增：原子更新三对象
	UpdateAllInTransaction(b *orderentity.Booking, r *orderentity.PickupRequest, o *orderentity.DriverOffer) error
}
