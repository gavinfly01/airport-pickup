package redisstore

import (
	"context"
	"encoding/json"
	"fmt"

	redis "github.com/redis/go-redis/v9"
)

// Client 使用 go-redis 实现
type Client struct{ cli *redis.Client }

// Options 与占位实现保持一致
type Options struct {
	Addr     string
	Password string
	DB       int
}

func New(opt Options) *Client {
	addr := opt.Addr
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	c := redis.NewClient(&redis.Options{Addr: addr, Password: opt.Password, DB: opt.DB})
	return &Client{cli: c}
}

func (c *Client) Ping(ctx context.Context) error {
	return c.cli.Ping(ctx).Err()
}

func (c *Client) AddPickupRequest(ctx context.Context, airport, vehicle string, req any, maxPrice float64) error {
	key := fmt.Sprintf("orderbook:requests:%s:%s", airport, vehicle)
	b, _ := json.Marshal(req)
	z := &redis.Z{Score: maxPrice, Member: b}
	return c.cli.ZAdd(ctx, key, *z).Err()
}

func (c *Client) AddDriverOffer(ctx context.Context, airport, vehicle string, offer any, price float64) error {
	key := fmt.Sprintf("orderbook:offers:%s:%s", airport, vehicle)
	b, _ := json.Marshal(offer)
	z := &redis.Z{Score: price, Member: b}
	return c.cli.ZAdd(ctx, key, *z).Err()
}

// RemovePickupRequest 删除 ZSET 中匹配 requestID 的成员
func (c *Client) RemovePickupRequest(ctx context.Context, airport, vehicle, requestID string) error {
	if requestID == "" {
		return fmt.Errorf("requestID is empty")
	}
	key := fmt.Sprintf("orderbook:requests:%s:%s", airport, vehicle)
	vals, err := c.cli.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		return err
	}
	for _, m := range vals {
		var obj map[string]any
		if json.Unmarshal([]byte(m), &obj) == nil {
			if idv, ok := obj["RequestID"]; ok && fmt.Sprint(idv) == requestID {
				// 按原始成员删除
				if err := c.cli.ZRem(ctx, key, m).Err(); err != nil {
					return err
				}
				break
			}
			if idv, ok := obj["ID"]; ok && fmt.Sprint(idv) == requestID { // 兼容实体 JSON
				if err := c.cli.ZRem(ctx, key, m).Err(); err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

// RemoveDriverOffer 删除 ZSET 中匹配 offerID 的成员
func (c *Client) RemoveDriverOffer(ctx context.Context, airport, vehicle, offerID string) error {
	if offerID == "" {
		return nil
	}
	key := fmt.Sprintf("orderbook:offers:%s:%s", airport, vehicle)
	vals, err := c.cli.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		return err
	}
	for _, m := range vals {
		var obj map[string]any
		if json.Unmarshal([]byte(m), &obj) == nil {
			if idv, ok := obj["OfferID"]; ok && fmt.Sprint(idv) == offerID {
				if err := c.cli.ZRem(ctx, key, m).Err(); err != nil {
					return err
				}
				break
			}
			if idv, ok := obj["ID"]; ok && fmt.Sprint(idv) == offerID { // 兼容实体 JSON
				if err := c.cli.ZRem(ctx, key, m).Err(); err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

func (c *Client) Close() error { return c.cli.Close() }
