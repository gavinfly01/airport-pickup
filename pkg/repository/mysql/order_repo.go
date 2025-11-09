package mysqlrepo

import (
	order "github.com/gavin/airport-pickup/internal/domain/order"
	orderentity "github.com/gavin/airport-pickup/internal/domain/order/entity"
	"gorm.io/gorm"
	"time"
)

type OrderRepository struct{ db *gorm.DB }

func NewOrderRepository(db *gorm.DB) order.OrderRepository { return &OrderRepository{db: db} }

// PickupRequest
func (r *OrderRepository) SavePickupRequest(p *orderentity.PickupRequest) error {
	m := &PickupRequest{
		ID: p.ID, PassengerID: p.PassengerID, AirportCode: p.AirportCode, VehicleType: p.VehicleType,
		DesiredTime: p.DesiredTime, MaxPricePerKm: p.MaxPricePerKm, PreferHighRating: p.PreferHighRating, Status: p.Status,
	}
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return r.db.Save(m).Error
}

func (r *OrderRepository) GetPickupRequestByID(id string) (*orderentity.PickupRequest, error) {
	var m PickupRequest
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &orderentity.PickupRequest{
		ID: m.ID, PassengerID: m.PassengerID, AirportCode: m.AirportCode, VehicleType: m.VehicleType,
		DesiredTime: m.DesiredTime, MaxPricePerKm: m.MaxPricePerKm, PreferHighRating: m.PreferHighRating, Status: m.Status,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}, nil
}

func (r *OrderRepository) ListPickupRequests() ([]*orderentity.PickupRequest, error) {
	var ms []PickupRequest
	if err := r.db.Find(&ms).Error; err != nil {
		return nil, err
	}
	res := make([]*orderentity.PickupRequest, 0, len(ms))
	for _, m := range ms {
		res = append(res, &orderentity.PickupRequest{
			ID: m.ID, PassengerID: m.PassengerID, AirportCode: m.AirportCode, VehicleType: m.VehicleType,
			DesiredTime: m.DesiredTime, MaxPricePerKm: m.MaxPricePerKm, PreferHighRating: m.PreferHighRating, Status: m.Status,
			CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
		})
	}
	return res, nil
}

func (r *OrderRepository) UpdatePickupRequest(p *orderentity.PickupRequest) error {
	return r.SavePickupRequest(p)
}

// DriverOffer
func (r *OrderRepository) SaveDriverOffer(o *orderentity.DriverOffer) error {
	m := &DriverOffer{
		ID: o.ID, DriverID: o.DriverID, AirportCode: o.AirportCode, VehicleType: o.VehicleType,
		AvailableFrom: o.AvailableFrom, AvailableTo: o.AvailableTo, PricePerKm: o.PricePerKm,
		Rating: o.Rating, Status: o.Status,
	}
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return r.db.Save(m).Error
}

func (r *OrderRepository) GetDriverOfferByID(id string) (*orderentity.DriverOffer, error) {
	var m DriverOffer
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &orderentity.DriverOffer{
		ID: m.ID, DriverID: m.DriverID, AirportCode: m.AirportCode, VehicleType: m.VehicleType,
		AvailableFrom: m.AvailableFrom, AvailableTo: m.AvailableTo, PricePerKm: m.PricePerKm,
		Rating: m.Rating, Status: m.Status,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}, nil
}

func (r *OrderRepository) ListDriverOffers() ([]*orderentity.DriverOffer, error) {
	var ms []DriverOffer
	if err := r.db.Find(&ms).Error; err != nil {
		return nil, err
	}
	res := make([]*orderentity.DriverOffer, 0, len(ms))
	for _, m := range ms {
		res = append(res, &orderentity.DriverOffer{
			ID: m.ID, DriverID: m.DriverID, AirportCode: m.AirportCode, VehicleType: m.VehicleType,
			AvailableFrom: m.AvailableFrom, AvailableTo: m.AvailableTo, PricePerKm: m.PricePerKm,
			Rating: m.Rating, Status: m.Status,
			CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
		})
	}
	return res, nil
}

func (r *OrderRepository) UpdateDriverOffer(o *orderentity.DriverOffer) error {
	return r.SaveDriverOffer(o)
}

// Booking
func (r *OrderRepository) SaveBooking(b *orderentity.Booking) error {
	m := &Booking{
		ID: b.ID, RequestID: b.RequestID, OfferID: b.OfferID, PassengerID: b.PassengerID, DriverID: b.DriverID,
		PricePerKm: b.PricePerKm, PlatformMarginPerKm: b.PlatformMarginPerKm, Status: b.Status,
	}
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return r.db.Save(m).Error
}

func (r *OrderRepository) GetBookingByID(id string) (*orderentity.Booking, error) {
	var m Booking
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &orderentity.Booking{
		ID: m.ID, RequestID: m.RequestID, OfferID: m.OfferID, PassengerID: m.PassengerID, DriverID: m.DriverID,
		PricePerKm: m.PricePerKm, PlatformMarginPerKm: m.PlatformMarginPerKm, Status: m.Status,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}, nil
}

func (r *OrderRepository) ListBookings() ([]*orderentity.Booking, error) {
	var ms []Booking
	if err := r.db.Find(&ms).Error; err != nil {
		return nil, err
	}
	res := make([]*orderentity.Booking, 0, len(ms))
	for _, m := range ms {
		res = append(res, &orderentity.Booking{
			ID: m.ID, RequestID: m.RequestID, OfferID: m.OfferID, PassengerID: m.PassengerID, DriverID: m.DriverID,
			PricePerKm: m.PricePerKm, PlatformMarginPerKm: m.PlatformMarginPerKm, Status: m.Status,
			CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
		})
	}
	return res, nil
}

func (r *OrderRepository) UpdateBooking(b *orderentity.Booking) error { return r.SaveBooking(b) }

func (r *OrderRepository) HasOngoingPickupRequest(passengerID string) (bool, error) {
	var cnt int64
	err := r.db.Model(&PickupRequest{}).
		Where("passenger_id = ? AND status IN ?", passengerID, []string{"open", "matched"}).
		Count(&cnt).Error
	return cnt > 0, err
}

func (r *OrderRepository) HasOngoingDriverOffer(driverID string) (bool, error) {
	var cnt int64
	err := r.db.Model(&DriverOffer{}).
		Where("driver_id = ? AND status IN ?", driverID, []string{"open", "matched"}).
		Count(&cnt).Error
	return cnt > 0, err
}

func (r *OrderRepository) UpdateAllInTransaction(b *orderentity.Booking, req *orderentity.PickupRequest, ofr *orderentity.DriverOffer) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		if b != nil {
			createdAt := b.CreatedAt
			if createdAt.IsZero() {
				createdAt = now
			}
			mB := &Booking{
				ID: b.ID, RequestID: b.RequestID, OfferID: b.OfferID, PassengerID: b.PassengerID, DriverID: b.DriverID,
				PricePerKm: b.PricePerKm, PlatformMarginPerKm: b.PlatformMarginPerKm, Status: b.Status, CreatedAt: createdAt,
			}
			mB.UpdatedAt = now
			if err := tx.Save(mB).Error; err != nil {
				return err
			}
		}
		if req != nil {
			createdAt := req.CreatedAt
			if createdAt.IsZero() {
				createdAt = now
			}
			mReq := &PickupRequest{
				ID: req.ID, PassengerID: req.PassengerID, AirportCode: req.AirportCode, VehicleType: req.VehicleType,
				DesiredTime: req.DesiredTime, MaxPricePerKm: req.MaxPricePerKm, PreferHighRating: req.PreferHighRating,
				Status: req.Status, CreatedAt: createdAt,
			}
			mReq.UpdatedAt = now
			if err := tx.Save(mReq).Error; err != nil {
				return err
			}
		}
		if ofr != nil {
			createdAt := ofr.CreatedAt
			if createdAt.IsZero() {
				createdAt = now
			}
			mOfr := &DriverOffer{
				ID: ofr.ID, DriverID: ofr.DriverID, AirportCode: ofr.AirportCode, VehicleType: ofr.VehicleType,
				AvailableFrom: ofr.AvailableFrom, AvailableTo: ofr.AvailableTo, PricePerKm: ofr.PricePerKm,
				Rating: ofr.Rating, Status: ofr.Status, CreatedAt: createdAt,
			}
			mOfr.UpdatedAt = now
			if err := tx.Save(mOfr).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
