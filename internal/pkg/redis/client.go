package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/maxfeizi04-cloud/recruitment-platform/internal/config"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func NewClient(cfg config.RedisConfig) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &Client{Client: rdb}, nil
}

func (c *Client) SetVerificationCode(ctx context.Context, phone, code string) error {
	key := fmt.Sprintf("sms:code:%s", phone)
	return c.Set(ctx, key, code, 2*time.Minute).Err()
}

func (c *Client) GetAndDeleteVerificationCode(ctx context.Context, phone, code string) (bool, error) {
	key := fmt.Sprintf("sms:code:%s", phone)
	stored, err := c.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("redis get: %w", err)
	}
	if stored == code {
		c.Del(ctx, key)
		return true, nil
	}
	return false, nil
}

func (c *Client) CanSendCode(ctx context.Context, phone string) (bool, error) {
	key := fmt.Sprintf("sms:limit:%s", phone)
	exists, err := c.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists: %w", err)
	}
	if exists > 0 {
		return false, nil
	}
	c.Set(ctx, key, "1", 60*time.Second)
	return true, nil
}
