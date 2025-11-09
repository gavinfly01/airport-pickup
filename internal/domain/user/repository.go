package user

import "github.com/gavin/airport-pickup/internal/domain/user/entity"

type PassengerRepository interface {
	Save(p *entity.Passenger) error
	GetByID(id string) (*entity.Passenger, error)
}

type DriverRepository interface {
	Save(d *entity.Driver) error
	GetByID(id string) (*entity.Driver, error)
}
