package mysqlrepo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDB_EmptyDSN(t *testing.T) {
	db, err := NewDB("")
	assert.Nil(t, db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty DSN")
}

func TestNewDB_InvalidDSN(t *testing.T) {
	db, err := NewDB("invalid_dsn")
	assert.Nil(t, db)
	assert.Error(t, err)
}
