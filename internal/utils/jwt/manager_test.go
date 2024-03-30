package jwt

import (
	"fmt"
	"github.com/mistandok/auth/internal/config"
	"github.com/mistandok/auth/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestManager_SuccessAccessTokenGeneration(t *testing.T) {
	manager := NewManager(&config.JWTConfig{
		JWTSecretKey:                 "secret",
		JWTAccessTokenExpireThrough:  1 * time.Minute,
		JWTRefreshTokenExpireThrough: 60 * time.Minute,
	})
	user := *testUser()
	accessToken, err := manager.GenerateAccessToken(user)
	require.NoError(t, err)

	_, err = manager.VerifyToken(accessToken)
	require.NoError(t, err)
}

func TestManager_SuccessRefreshTokenGeneration(t *testing.T) {
	manager := NewManager(&config.JWTConfig{
		JWTSecretKey:                 "secret",
		JWTAccessTokenExpireThrough:  1 * time.Minute,
		JWTRefreshTokenExpireThrough: 60 * time.Minute,
	})
	user := *testUser()
	accessToken, err := manager.GenerateRefreshToken(user)
	require.NoError(t, err)

	_, err = manager.VerifyToken(accessToken)
	require.NoError(t, err)
}

func TestManager_SuccessAccessTokenGenerationWithClaimCheck(t *testing.T) {
	manager := NewManager(&config.JWTConfig{
		JWTSecretKey:                 "secret",
		JWTAccessTokenExpireThrough:  1 * time.Minute,
		JWTRefreshTokenExpireThrough: 60 * time.Minute,
	})
	user := *testUser()
	accessToken, err := manager.GenerateAccessToken(user)
	require.NoError(t, err)

	userClaims, err := manager.VerifyToken(accessToken)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%d", user.ID), userClaims.UserID)
	require.Equal(t, user.Name, userClaims.UserName)
	require.Equal(t, string(user.Role), userClaims.Role)
}

func TestManager_AccessTokenExpired(t *testing.T) {
	t.Parallel()
	manager := NewManager(&config.JWTConfig{
		JWTSecretKey:                 "secret",
		JWTAccessTokenExpireThrough:  2 * time.Second,
		JWTRefreshTokenExpireThrough: 60 * time.Minute,
	})
	user := *testUser()
	accessToken, err := manager.GenerateAccessToken(user)
	require.NoError(t, err)

	time.Sleep(3 * time.Second)

	_, err = manager.VerifyToken(accessToken)
	require.Error(t, err)
}

func TestManager_RefreshTokenExpired(t *testing.T) {
	t.Parallel()
	manager := NewManager(&config.JWTConfig{
		JWTSecretKey:                 "secret",
		JWTAccessTokenExpireThrough:  2 * time.Second,
		JWTRefreshTokenExpireThrough: 2 * time.Second,
	})
	user := *testUser()
	accessToken, err := manager.GenerateRefreshToken(user)
	require.NoError(t, err)

	time.Sleep(3 * time.Second)

	_, err = manager.VerifyToken(accessToken)
	require.Error(t, err)
}

func testUser() *model.User {
	return &model.User{
		ID:        1,
		Name:      "anton",
		Email:     "email@mail.com",
		Role:      "USER",
		Password:  "password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
