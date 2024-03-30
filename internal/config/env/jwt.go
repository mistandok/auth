package env

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/mistandok/auth/internal/config"
)

const (
	jwtKeyEnvName           = "JWT_SECRET_KEY"
	jwtAccessExpireEnvName  = "JWT_ACCESS_TOKEN_EXPIRES_MIN"
	jwtRefreshExpireEnvName = "JWT_REFRESH_TOKEN_EXPIRES_MIN"
)

// JWTConfigSearcher logger config searcher.
type JWTConfigSearcher struct{}

// NewJWTConfigSearcher get instance for logger config searcher.
func NewJWTConfigSearcher() *JWTConfigSearcher {
	return &JWTConfigSearcher{}
}

// Get config for password.
func (s *JWTConfigSearcher) Get() (*config.JWTConfig, error) {
	jwtSecret := os.Getenv(jwtKeyEnvName)
	if len(jwtSecret) == 0 {
		return nil, errors.New("не найден секрет для JWT")
	}

	accessExpireStr := os.Getenv(jwtAccessExpireEnvName)
	if len(accessExpireStr) == 0 {
		return nil, errors.New("не найдено время жизни access токена")
	}

	accessExpireInt, err := strconv.Atoi(accessExpireStr)
	if err != nil {
		return nil, errors.New("некорректное время жизни access токена")
	}

	refreshExpireStr := os.Getenv(jwtRefreshExpireEnvName)
	if len(accessExpireStr) == 0 {
		return nil, errors.New("не найдено время жизни access токена")
	}

	refreshExpireInt, err := strconv.Atoi(refreshExpireStr)
	if err != nil {
		return nil, errors.New("некорректное время жизни access токена")
	}

	return &config.JWTConfig{
		JWTSecretKey:                    jwtSecret,
		JWTAccessTokenExpireThroughMin:  time.Duration(accessExpireInt) * time.Minute,
		JWTRefreshTokenExpireThroughMin: time.Duration(refreshExpireInt) * time.Minute,
	}, nil
}
