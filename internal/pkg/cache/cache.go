package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client Redis 缓存客户端
type Client struct {
	rdb *redis.Client
}

// New 创建缓存客户端
func New(rdb *redis.Client) *Client {
	return &Client{rdb: rdb}
}

// Get 从缓存读取并反序列化到 dest
func (c *Client) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	data, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("cache get: %w", err)
	}
	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return false, fmt.Errorf("cache unmarshal: %w", err)
	}
	return true, nil
}

// Set 序列化 value 并写入缓存
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal: %w", err)
	}
	return c.rdb.Set(ctx, key, data, ttl).Err()
}

// Delete 删除缓存 key
func (c *Client) Delete(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

// DeletePattern 删除匹配 pattern 的所有 key（如 "jobs:*"）
func (c *Client) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.rdb.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		if err := c.rdb.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// ── 常用 TTL ──

const (
	TTLShort  = 1 * time.Minute   // 1分钟
	TTLMedium = 5 * time.Minute   // 5分钟
	TTLLong   = 30 * time.Minute  // 30分钟
	TTLHour   = 1 * time.Hour     // 1小时
)
