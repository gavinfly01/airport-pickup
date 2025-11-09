package payments

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWalletClient_Charge_Success(t *testing.T) {
	client := NewWalletClient()
	err := client.Charge("booking123", 1000)
	assert.NoError(t, err)
}

func TestWalletClient_Charge_InvalidAmount(t *testing.T) {
	client := NewWalletClient()
	err := client.Charge("booking123", -100)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid amount")
}
