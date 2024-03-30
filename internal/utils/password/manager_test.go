package password

import (
	"testing"

	"github.com/mistandok/auth/internal/config"
	"github.com/stretchr/testify/require"
)

func TestManager_SuccessPassCompare(t *testing.T) {
	manager := NewManager(&config.PasswordConfig{PasswordSalt: "test_salt"})

	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "first case",
			password: "best_password",
		},
		{
			name:     "second case",
			password: "amazing_creator",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hash, err := manager.HashPassword(test.password)
			require.NoError(t, err)

			isEqual := manager.CheckPasswordHash(test.password, hash)
			require.Equal(t, isEqual, true)
		})
	}
}

func TestManager_FailPassCompare(t *testing.T) {
	manager := NewManager(&config.PasswordConfig{PasswordSalt: "test_salt"})

	tests := []struct {
		name         string
		password     string
		fakePassword string
	}{
		{
			name:         "first case",
			password:     "best_password",
			fakePassword: "test_fake",
		},
		{
			name:         "second case",
			password:     "amazing_creator",
			fakePassword: "test_fake",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hash, err := manager.HashPassword(test.password)
			require.NoError(t, err)

			isEqual := manager.CheckPasswordHash(test.fakePassword, hash)
			require.Equal(t, isEqual, false)
		})
	}
}
