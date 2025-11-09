package mysqlrepo

import (
	"errors"
	"time"

	user "github.com/gavin/airport-pickup/internal/domain/user"
	userentity "github.com/gavin/airport-pickup/internal/domain/user/entity"
	"gorm.io/gorm"
)

type PassengerRepository struct{ db *gorm.DB }

type DriverRepository struct{ db *gorm.DB }

func NewPassengerRepository(db *gorm.DB) user.PassengerRepository {
	return &PassengerRepository{db: db}
}
func NewDriverRepository(db *gorm.DB) user.DriverRepository { return &DriverRepository{db: db} }

// Passenger
func (r *PassengerRepository) Save(p *userentity.Passenger) error {
	if p == nil || p.ID == "" {
		return errors.New("invalid passenger")
	}
	m := &Passenger{ID: p.ID, Name: p.Name}
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return r.db.Save(m).Error
}

func (r *PassengerRepository) GetByID(id string) (*userentity.Passenger, error) {
	var m Passenger
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &userentity.Passenger{ID: m.ID, Name: m.Name, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

// Driver
func (r *DriverRepository) Save(d *userentity.Driver) error {
	if d == nil || d.ID == "" {
		return errors.New("invalid driver")
	}
	m := &Driver{ID: d.ID, Name: d.Name, Rating: d.Rating}
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return r.db.Save(m).Error
}

func (r *DriverRepository) GetByID(id string) (*userentity.Driver, error) {
	var m Driver
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &userentity.Driver{ID: m.ID, Name: m.Name, Rating: m.Rating, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}
