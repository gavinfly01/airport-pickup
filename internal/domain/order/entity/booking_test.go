package entity

import (
	"testing"
)

func TestBooking_MarkCompleted(t *testing.T) {
	b := &Booking{Status: "created"}
	err := b.MarkCompleted()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if b.Status != "completed" {
		t.Errorf("expected status 'completed', got %v", b.Status)
	}

	b2 := &Booking{Status: "cancelled"}
	err2 := b2.MarkCompleted()
	if err2 == nil {
		t.Errorf("expected error for invalid status, got nil")
	}
}
