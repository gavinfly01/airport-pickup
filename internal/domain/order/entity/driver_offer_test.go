package entity

import (
	"testing"
)

func TestDriverOffer_MarkMatched(t *testing.T) {
	o := &DriverOffer{Status: "open"}
	err := o.MarkMatched()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if o.Status != "matched" {
		t.Errorf("expected status 'matched', got %v", o.Status)
	}

	o2 := &DriverOffer{Status: "cancelled"}
	err2 := o2.MarkMatched()
	if err2 == nil {
		t.Errorf("expected error for invalid status, got nil")
	}
}

func TestDriverOffer_MarkCompleted(t *testing.T) {
	o := &DriverOffer{Status: "matched"}
	err := o.MarkCompleted()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if o.Status != "completed" {
		t.Errorf("expected status 'completed', got %v", o.Status)
	}

	o2 := &DriverOffer{Status: "open"}
	err2 := o2.MarkCompleted()
	if err2 == nil {
		t.Errorf("expected error for invalid status, got nil")
	}
}
