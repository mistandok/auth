package rs

import (
	"github.com/gomodule/redigo/redis"
	"github.com/mistandok/auth/internal/temp_redis"
)

var _ temp_redis.Client = (*redisClient)(nil)

type redisClient struct {
	masterDB temp_redis.DB
}

func New(pool *redis.Pool) temp_redis.Client {
	return &redisClient{masterDB: NewRs(pool)}
}

func (c *redisClient) DB() temp_redis.DB {
	return c.masterDB
}

func (c *redisClient) Close() error {
	if c.masterDB != nil {
		return c.masterDB.Close()
	}

	return nil
}
