package service

import (
	"errors"
	settlemententity "github.com/gavin/airport-pickup/internal/domain/settlement/entity"
)

// PaymentService defines interaction with external wallet/payment gateway.
type PaymentService interface {
	Charge(bookingID string, amountCents int64) error
}

type CreatePaymentTransactionCmd struct {
	BookingID   string
	AmountCents int64
	Status      string
}

type PaymentTransactionService struct{}

func NewPaymentTransactionService() *PaymentTransactionService {
	return &PaymentTransactionService{}
}

func (s *PaymentTransactionService) CreatePaymentTransaction(cmd *CreatePaymentTransactionCmd) (*settlemententity.PaymentTransaction, error) {
	if cmd.BookingID == "" {
		return nil, errors.New("booking_id required")
	}
	if cmd.AmountCents < 0 {
		return nil, errors.New("amount_cents must be >= 0")
	}
	if cmd.Status == "" {
		return nil, errors.New("status required")
	}
	return &settlemententity.PaymentTransaction{
		ID:          "",
		BookingID:   cmd.BookingID,
		AmountCents: cmd.AmountCents,
		Status:      cmd.Status,
	}, nil
}
