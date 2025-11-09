package app

import (
	"errors"
	"github.com/gavin/airport-pickup/pkg/util"

	evt "github.com/gavin/airport-pickup/internal/domain/eventbus"
	order "github.com/gavin/airport-pickup/internal/domain/order"
	settlement "github.com/gavin/airport-pickup/internal/domain/settlement"
	settlesvc "github.com/gavin/airport-pickup/internal/domain/settlement/service"
)

type SettlementAppService struct {
	repo      settlement.SettlementRepository
	orderRepo order.OrderRepository
	pay       settlesvc.PaymentService
	bus       evt.EventBus

	paymentTxService  *settlesvc.PaymentTransactionService
	settlementService settlesvc.SettlementService
}

func NewSettlementAppService(repo settlement.SettlementRepository, orderRepo order.OrderRepository, pay settlesvc.PaymentService, bus evt.EventBus) *SettlementAppService {
	return &SettlementAppService{
		repo:              repo,
		orderRepo:         orderRepo,
		pay:               pay,
		bus:               bus,
		paymentTxService:  settlesvc.NewPaymentTransactionService(),
		settlementService: settlesvc.NewSettlementService(),
	}
}

func (s *SettlementAppService) TriggerPayment(bookingID string) error {
	return s.OnOrderCompleted(bookingID)
}

// OnOrderCompleted orchestrates payment -> settlement -> revenue update
func (s *SettlementAppService) OnOrderCompleted(bookingID string) error {
	b, err := s.orderRepo.GetBookingByID(bookingID)
	if err != nil || b == nil {
		return errors.New("booking not found")
	}
	// naive amount calculation: assume 10km for demo only
	amountCents := int64(b.PricePerKm * 100.0 * 10.0)
	platformRevenueCents := int64(b.PlatformMarginPerKm * 100.0 * 10.0)
	if amountCents < 0 {
		amountCents = 0
	}
	if platformRevenueCents < 0 {
		platformRevenueCents = 0
	}
	// charge
	if err := s.pay.Charge(bookingID, amountCents); err != nil {
		return err
	}
	// persist payment (通过领域服务)
	ptx, err := s.paymentTxService.CreatePaymentTransaction(&settlesvc.CreatePaymentTransactionCmd{
		BookingID:   bookingID,
		AmountCents: amountCents,
		Status:      "success",
	})
	if err != nil {
		return err
	}
	// settlement record (通过领域服务)
	sr, err := s.settlementService.CreateSettlementRecord(&settlesvc.CreateSettlementRecordCmd{
		BookingID:            bookingID,
		DriverID:             b.DriverID,
		PassengerID:          b.PassengerID,
		AmountCents:          amountCents,
		PlatformRevenueCents: platformRevenueCents,
	})
	if err != nil {
		return err
	}
	// revenue record (通过领域服务)
	rr, err := s.settlementService.CreateRevenueRecord(&settlesvc.CreateRevenueRecordCmd{
		BookingID:  bookingID,
		DeltaCents: platformRevenueCents,
	})
	if err != nil {
		return err
	}
	ptx.ID = util.NewID()
	sr.ID = util.NewID()
	rr.ID = util.NewID()
	if err := s.repo.SaveAllInTransaction(ptx, sr, rr); err != nil {
		return err
	}
	// publish events
	s.bus.Publish(evt.PaymentSucceeded{BookingID: bookingID, AmountCents: amountCents})
	s.bus.Publish(evt.SettlementCreated{BookingID: bookingID})
	s.bus.Publish(evt.RevenueUpdated{BookingID: bookingID, DeltaCents: platformRevenueCents})
	return nil
}
