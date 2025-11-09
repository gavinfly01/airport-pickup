package mysqlrepo

import (
	orderentity "github.com/gavin/airport-pickup/internal/domain/order/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

func newTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&PickupRequest{})
	return db
}

func TestSaveAndGetPickupRequest(t *testing.T) {
	db := newTestDB()
	repo := NewOrderRepository(db)
	pr := &orderentity.PickupRequest{
		ID: "r1", PassengerID: "p1", AirportCode: "PVG", VehicleType: "Sedan", DesiredTime: time.Now(), MaxPricePerKm: 10, PreferHighRating: true, Status: "open",
	}
	err := repo.SavePickupRequest(pr)
	assert.NoError(t, err)

	got, err := repo.GetPickupRequestByID("r1")
	assert.NoError(t, err)
	assert.Equal(t, "r1", got.ID)
	assert.Equal(t, "p1", got.PassengerID)
}

func TestGetPickupRequestByID_NotFound(t *testing.T) {
	db := newTestDB()
	repo := NewOrderRepository(db)
	_, err := repo.GetPickupRequestByID("not_exist")
	assert.Error(t, err)
}

func TestListPickupRequests(t *testing.T) {
	db := newTestDB()
	repo := NewOrderRepository(db)
	pr := &orderentity.PickupRequest{ID: "r2", PassengerID: "p2", AirportCode: "SHA", VehicleType: "SUV", DesiredTime: time.Now(), MaxPricePerKm: 20, PreferHighRating: false, Status: "open"}
	repo.SavePickupRequest(pr)
	list, err := repo.ListPickupRequests()
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "r2", list[0].ID)
}
