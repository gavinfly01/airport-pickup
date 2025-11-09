package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewID_Length(t *testing.T) {
	id := NewID()
	assert.Equal(t, 32, len(id))
}

func TestNewID_Unique(t *testing.T) {
	ids := make(map[string]struct{})
	for i := 0; i < 1000; i++ {
		id := NewID()
		_, exists := ids[id]
		assert.False(t, exists, "duplicate id found")
		ids[id] = struct{}{}
	}
}
