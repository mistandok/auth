package env

import (
	"errors"
	"os"

	"github.com/mistandok/auth/internal/config"
)

const (
	whiteListHostEnv = "WHITE_LIST_REDIS_HOST"
	whiteListPortEnv = "WHITE_LIST_REDIS_PORT"
)

// WhiteListRedisCfgSearcher searcher for white list redis config.
type WhiteListRedisCfgSearcher struct{}

// NewWhiteListRedisCfgSearcher get instance searcher for white list redis config.
func NewWhiteListRedisCfgSearcher() *WhiteListRedisCfgSearcher {
	return &WhiteListRedisCfgSearcher{}
}

// Get config for PG connection.
func (s *WhiteListRedisCfgSearcher) Get() (*config.WhiteListRedisConfig, error) {
	dbHost := os.Getenv(whiteListHostEnv)
	if len(dbHost) == 0 {
		return nil, errors.New("white list redis host not found")
	}

	dbPort := os.Getenv(whiteListPortEnv)
	if len(dbPort) == 0 {
		return nil, errors.New("white list redis port not found")
	}

	return &config.WhiteListRedisConfig{
		Host: dbHost,
		Port: dbPort,
	}, nil
}
