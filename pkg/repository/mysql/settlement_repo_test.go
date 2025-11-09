package mysqlrepo

import (
	settlemententity "github.com/gavin/airport-pickup/internal/domain/settlement/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

func newTestDBSettlement() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&PaymentTransaction{}, &SettlementRecord{}, &RevenueRecord{})
	return db
}

func TestSaveAndGetPaymentTransaction(t *testing.T) {
	db := newTestDBSettlement()
	repo := NewSettlementRepository(db)
	pt := &settlemententity.PaymentTransaction{ID: "pt1", BookingID: "b1", AmountCents: 1000, Status: "paid", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err := repo.SavePaymentTransaction(pt)
	assert.NoError(t, err)
	got, err := repo.GetPaymentTransactionByID("pt1")
	assert.NoError(t, err)
	assert.Equal(t, "pt1", got.ID)
	assert.Equal(t, "paid", got.Status)
}

func TestGetPaymentTransactionByID_NotFound(t *testing.T) {
	db := newTestDBSettlement()
	repo := NewSettlementRepository(db)
	_, err := repo.GetPaymentTransactionByID("not_exist")
	assert.Error(t, err)
}

func TestSaveAndGetSettlementRecord(t *testing.T) {
	db := newTestDBSettlement()
	repo := NewSettlementRepository(db)
	sr := &settlemententity.SettlementRecord{ID: "sr1", BookingID: "b2", DriverID: "d1", PassengerID: "p1", AmountCents: 2000, PlatformRevenueCents: 100, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err := repo.SaveSettlementRecord(sr)
	assert.NoError(t, err)
	got, err := repo.GetSettlementRecordByID("sr1")
	assert.NoError(t, err)
	assert.Equal(t, "sr1", got.ID)
	assert.Equal(t, int64(2000), got.AmountCents)
}

func TestSaveRevenueRecord(t *testing.T) {
	db := newTestDBSettlement()
	repo := NewSettlementRepository(db)
	rr := &settlemententity.RevenueRecord{ID: "rev1", BookingID: "b3", DeltaCents: 50, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err := repo.SaveRevenueRecord(rr)
	assert.NoError(t, err)
}
