package env

import (
	"errors"
	"github.com/mistandok/auth/internal/config"
	"os"
)

const (
	passwordSalt = "PASSWORD_SALT"
)

// PasswordConfigSearcher logger config searcher.
type PasswordConfigSearcher struct{}

// NewPasswordConfigSearcher get instance for passford config searcher.
func NewPasswordConfigSearcher() *PasswordConfigSearcher {
	return &PasswordConfigSearcher{}
}

// Get config for password.
func (s *PasswordConfigSearcher) Get() (*config.PasswordConfig, error) {
	salt := os.Getenv(passwordSalt)
	if len(salt) == 0 {
		return nil, errors.New("не найдена соль для пароля")
	}

	return &config.PasswordConfig{
		PasswordSalt: salt,
	}, nil
}
