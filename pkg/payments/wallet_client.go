package payments

import "fmt"

// WalletClient is a simple stub implementing the PaymentService by charging successfully.
type WalletClient struct{}

func NewWalletClient() *WalletClient { return &WalletClient{} }

func (w *WalletClient) Charge(bookingID string, amountCents int64) error {
	// In real world, call external gateway; here succeed if amount >= 0
	if amountCents < 0 {
		return fmt.Errorf("invalid amount")
	}
	return nil
}
