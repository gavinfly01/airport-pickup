package redisstore

import (
	"context"
	"testing"

	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func newMockClient() *Client {
	cli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // 测试库
	})
	return &Client{cli: cli}
}

func TestPing(t *testing.T) {
	c := newMockClient()
	err := c.Ping(context.Background())
	assert.NoError(t, err)
}

func TestAddPickupRequest(t *testing.T) {
	c := newMockClient()
	err := c.AddPickupRequest(context.Background(), "PVG", "Sedan", map[string]string{"id": "req1"}, 100)
	assert.NoError(t, err)
}

func TestAddDriverOffer(t *testing.T) {
	c := newMockClient()
	err := c.AddDriverOffer(context.Background(), "PVG", "Sedan", map[string]string{"id": "offer1"}, 80)
	assert.NoError(t, err)
}

func TestRemovePickupRequest_EmptyID(t *testing.T) {
	c := newMockClient()
	err := c.RemovePickupRequest(context.Background(), "PVG", "Sedan", "")
	assert.Error(t, err)
}
