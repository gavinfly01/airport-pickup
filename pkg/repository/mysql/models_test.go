package mysqlrepo

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPassengerInit(t *testing.T) {
	p := Passenger{ID: "p1", Name: "Alice", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	assert.Equal(t, "p1", p.ID)
	assert.Equal(t, "Alice", p.Name)
}

func TestDriverInit(t *testing.T) {
	d := Driver{ID: "d1", Name: "Bob", Rating: 4.8, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	assert.Equal(t, "d1", d.ID)
	assert.Equal(t, "Bob", d.Name)
	assert.Equal(t, 4.8, d.Rating)
}

func TestPickupRequestInit(t *testing.T) {
	r := PickupRequest{ID: "r1", PassengerID: "p1", AirportCode: "PVG", VehicleType: "Sedan", DesiredTime: time.Now(), MaxPricePerKm: 10, PreferHighRating: true, Status: "open", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	assert.Equal(t, "r1", r.ID)
	assert.Equal(t, "p1", r.PassengerID)
	assert.Equal(t, "PVG", r.AirportCode)
	assert.Equal(t, "Sedan", r.VehicleType)
	assert.Equal(t, "open", r.Status)
}

func TestDriverOfferInit(t *testing.T) {
	do := DriverOffer{ID: "o1", DriverID: "d1", AirportCode: "PVG", VehicleType: "SUV", AvailableFrom: time.Now(), AvailableTo: time.Now(), PricePerKm: 12, Rating: 4.9, Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	assert.Equal(t, "o1", do.ID)
	assert.Equal(t, "d1", do.DriverID)
	assert.Equal(t, "SUV", do.VehicleType)
	assert.Equal(t, "active", do.Status)
}
