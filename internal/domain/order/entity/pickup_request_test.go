package entity

import (
	"testing"
)

func TestPickupRequest_MarkMatched(t *testing.T) {
	r := &PickupRequest{Status: "open"}
	err := r.MarkMatched()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if r.Status != "matched" {
		t.Errorf("expected status 'matched', got %v", r.Status)
	}

	r2 := &PickupRequest{Status: "cancelled"}
	err2 := r2.MarkMatched()
	if err2 == nil {
		t.Errorf("expected error for invalid status, got nil")
	}
}

func TestPickupRequest_MarkCompleted(t *testing.T) {
	r := &PickupRequest{Status: "matched"}
	err := r.MarkCompleted()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if r.Status != "completed" {
		t.Errorf("expected status 'completed', got %v", r.Status)
	}

	r2 := &PickupRequest{Status: "open"}
	err2 := r2.MarkCompleted()
	if err2 == nil {
		t.Errorf("expected error for invalid status, got nil")
	}
}
