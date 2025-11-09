package settlement

import (
	settlemententity "github.com/gavin/airport-pickup/internal/domain/settlement/entity"
)

type SettlementRepository interface {
	SavePaymentTransaction(t *settlemententity.PaymentTransaction) error
	GetPaymentTransactionByID(id string) (*settlemententity.PaymentTransaction, error)

	SaveSettlementRecord(r *settlemententity.SettlementRecord) error
	GetSettlementRecordByID(id string) (*settlemententity.SettlementRecord, error)

	SaveRevenueRecord(r *settlemententity.RevenueRecord) error
	ListRevenueRecords() ([]*settlemententity.RevenueRecord, error)

	// 原子保存三对象
	SaveAllInTransaction(ptx *settlemententity.PaymentTransaction, sr *settlemententity.SettlementRecord, rr *settlemententity.RevenueRecord) error
}
