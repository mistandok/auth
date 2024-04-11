package env

import (
	"errors"
	"os"

	"github.com/mistandok/auth/internal/config"
)

const (
	redisHostEnv = "REDIS_HOST"
	redisPortEnv = "REDIS_PORT"
)

// RedisCfgSearcher searcher for white list redis config.
type RedisCfgSearcher struct{}

// NewRedisCfgSearcher get instance searcher for white list redis config.
func NewRedisCfgSearcher() *RedisCfgSearcher {
	return &RedisCfgSearcher{}
}

// Get config for PG connection.
func (s *RedisCfgSearcher) Get() (*config.RedisConfig, error) {
	dbHost := os.Getenv(redisHostEnv)
	if len(dbHost) == 0 {
		return nil, errors.New("white list redis host not found")
	}

	dbPort := os.Getenv(redisPortEnv)
	if len(dbPort) == 0 {
		return nil, errors.New("white list redis port not found")
	}

	return &config.RedisConfig{
		Host: dbHost,
		Port: dbPort,
	}, nil
}
