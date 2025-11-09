package mysqlrepo

import (
	settlement "github.com/gavin/airport-pickup/internal/domain/settlement"
	settlemententity "github.com/gavin/airport-pickup/internal/domain/settlement/entity"
	"gorm.io/gorm"
	"time"
)

type SettlementRepository struct{ db *gorm.DB }

func NewSettlementRepository(db *gorm.DB) settlement.SettlementRepository {
	return &SettlementRepository{db: db}
}

func (r *SettlementRepository) SavePaymentTransaction(t *settlemententity.PaymentTransaction) error {
	now := time.Now()
	m := &PaymentTransaction{ID: t.ID, BookingID: t.BookingID, AmountCents: t.AmountCents, Status: t.Status}
	m.CreatedAt = now
	m.UpdatedAt = now
	return r.db.Save(m).Error
}

func (r *SettlementRepository) GetPaymentTransactionByID(id string) (*settlemententity.PaymentTransaction, error) {
	var m PaymentTransaction
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &settlemententity.PaymentTransaction{ID: m.ID, BookingID: m.BookingID, AmountCents: m.AmountCents, Status: m.Status, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func (r *SettlementRepository) SaveSettlementRecord(s *settlemententity.SettlementRecord) error {
	now := time.Now()
	m := &SettlementRecord{ID: s.ID, BookingID: s.BookingID, DriverID: s.DriverID, PassengerID: s.PassengerID, AmountCents: s.AmountCents, PlatformRevenueCents: s.PlatformRevenueCents}
	m.CreatedAt = now
	m.UpdatedAt = now
	return r.db.Save(m).Error
}

func (r *SettlementRepository) GetSettlementRecordByID(id string) (*settlemententity.SettlementRecord, error) {
	var m SettlementRecord
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &settlemententity.SettlementRecord{ID: m.ID, BookingID: m.BookingID, DriverID: m.DriverID, PassengerID: m.PassengerID, AmountCents: m.AmountCents, PlatformRevenueCents: m.PlatformRevenueCents, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func (r *SettlementRepository) SaveRevenueRecord(rr *settlemententity.RevenueRecord) error {
	now := time.Now()
	m := &RevenueRecord{ID: rr.ID, BookingID: rr.BookingID, DeltaCents: rr.DeltaCents}
	m.CreatedAt = now
	m.UpdatedAt = now
	return r.db.Save(m).Error
}

func (r *SettlementRepository) ListRevenueRecords() ([]*settlemententity.RevenueRecord, error) {
	var ms []RevenueRecord
	if err := r.db.Find(&ms).Error; err != nil {
		return nil, err
	}
	res := make([]*settlemententity.RevenueRecord, 0, len(ms))
	for _, m := range ms {
		res = append(res, &settlemententity.RevenueRecord{ID: m.ID, BookingID: m.BookingID, DeltaCents: m.DeltaCents, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt})
	}
	return res, nil
}

func (r *SettlementRepository) SaveAllInTransaction(ptx *settlemententity.PaymentTransaction, sr *settlemententity.SettlementRecord, rr *settlemententity.RevenueRecord) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		mPTX := &PaymentTransaction{ID: ptx.ID, BookingID: ptx.BookingID, AmountCents: ptx.AmountCents, Status: ptx.Status}
		mPTX.CreatedAt = now
		mPTX.UpdatedAt = now
		if err := tx.Save(mPTX).Error; err != nil {
			return err
		}
		mSR := &SettlementRecord{ID: sr.ID, BookingID: sr.BookingID, DriverID: sr.DriverID, PassengerID: sr.PassengerID, AmountCents: sr.AmountCents, PlatformRevenueCents: sr.PlatformRevenueCents}
		mSR.CreatedAt = now
		mSR.UpdatedAt = now
		if err := tx.Save(mSR).Error; err != nil {
			return err
		}
		mRR := &RevenueRecord{ID: rr.ID, BookingID: rr.BookingID, DeltaCents: rr.DeltaCents}
		mRR.CreatedAt = now
		mRR.UpdatedAt = now
		if err := tx.Save(mRR).Error; err != nil {
			return err
		}
		return nil
	})
}
