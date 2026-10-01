package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func New(addr, password string, db int) (*Client, error) {
	if addr == "" {
		return nil, fmt.Errorf("redis address is required")
	}
	return &Client{rdb: redis.NewClient(&redis.Options{
		Addr: addr, Password: password, DB: db,
	})}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	if c == nil || c.rdb == nil {
		return fmt.Errorf("redis client is not configured")
	}
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return c.rdb.Ping(pingCtx).Err()
}

func (c *Client) Close() error {
	if c == nil || c.rdb == nil {
		return nil
	}
	return c.rdb.Close()
}
