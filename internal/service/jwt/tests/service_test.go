package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/mistandok/auth/internal/repository/mocks"
	autService "github.com/mistandok/auth/internal/service"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"

	"github.com/mistandok/auth/internal/config"
	"github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/internal/service/jwt"
	"github.com/stretchr/testify/require"
)

func TestService_SuccessAccessTokenGeneration(t *testing.T) {
	logger := zerolog.Nop()
	jwtConfig := &config.JWTConfig{
		JWTSecretKey:                 "secret",
		JWTAccessTokenExpireThrough:  1 * time.Minute,
		JWTRefreshTokenExpireThrough: 60 * time.Minute,
	}
	whiteListRepo := mocks.NewWhiteListRepository(t)

	service := jwt.NewService(&logger, jwtConfig, whiteListRepo)
	user := *testUser()
	accessToken, err := service.GenerateAccessToken(user)
	require.NoError(t, err)

	_, err = service.VerifyAccessToken(accessToken)
	require.NoError(t, err)
}

func TestService_SuccessRefreshTokenGeneration(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()
	jwtConfig := &config.JWTConfig{
		JWTSecretKey:                 "secret",
		JWTAccessTokenExpireThrough:  1 * time.Minute,
		JWTRefreshTokenExpireThrough: 60 * time.Minute,
	}
	user := *testUser()

	whiteListRepo := mocks.NewWhiteListRepository(t)
	whiteListRepo.On(
		"Set",
		ctx,
		user.ID,
		mock.AnythingOfType("string"),
		jwtConfig.JWTRefreshTokenExpireThrough,
	).Return(nil).Once()

	service := jwt.NewService(&logger, jwtConfig, whiteListRepo)
	refreshToken, err := service.GenerateRefreshToken(ctx, user)
	require.NoError(t, err)

	whiteListRepo.On(
		"Get",
		ctx,
		user.ID,
	).Return(refreshToken, nil).Once()

	_, err = service.VerifyRefreshToken(ctx, refreshToken)
	require.NoError(t, err)
}

func TestService_SuccessAccessTokenGenerationWithClaimCheck(t *testing.T) {
	logger := zerolog.Nop()
	jwtConfig := &config.JWTConfig{
		JWTSecretKey:                 "secret",
		JWTAccessTokenExpireThrough:  1 * time.Minute,
		JWTRefreshTokenExpireThrough: 60 * time.Minute,
	}
	whiteListRepo := mocks.NewWhiteListRepository(t)

	service := jwt.NewService(&logger, jwtConfig, whiteListRepo)
	user := *testUser()
	accessToken, err := service.GenerateAccessToken(user)
	require.NoError(t, err)

	userClaims, err := service.VerifyAccessToken(accessToken)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%d", user.ID), userClaims.UserID)
	require.Equal(t, user.Name, userClaims.UserName)
	require.Equal(t, string(user.Role), userClaims.Role)
}

func TestService_AccessTokenExpired(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	jwtConfig := &config.JWTConfig{
		JWTSecretKey:                 "secret",
		JWTAccessTokenExpireThrough:  1 * time.Second,
		JWTRefreshTokenExpireThrough: 60 * time.Minute,
	}
	whiteListRepo := mocks.NewWhiteListRepository(t)

	service := jwt.NewService(&logger, jwtConfig, whiteListRepo)
	user := *testUser()
	accessToken, err := service.GenerateAccessToken(user)
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	_, err = service.VerifyAccessToken(accessToken)
	require.Error(t, err)
}

func TestService_RefreshTokenExpired(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := zerolog.Nop()
	jwtConfig := &config.JWTConfig{
		JWTSecretKey:                 "secret",
		JWTAccessTokenExpireThrough:  1 * time.Minute,
		JWTRefreshTokenExpireThrough: 1 * time.Second,
	}
	user := *testUser()

	whiteListRepo := mocks.NewWhiteListRepository(t)
	whiteListRepo.On(
		"Set",
		ctx,
		user.ID,
		mock.AnythingOfType("string"),
		jwtConfig.JWTRefreshTokenExpireThrough,
	).Return(nil).Once()

	service := jwt.NewService(&logger, jwtConfig, whiteListRepo)

	refreshToken, err := service.GenerateRefreshToken(ctx, user)
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	_, err = service.VerifyRefreshToken(ctx, refreshToken)
	require.Error(t, err)
}

func TestService_RefreshTokenDiscredited(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()
	jwtConfig := &config.JWTConfig{
		JWTSecretKey:                 "secret",
		JWTAccessTokenExpireThrough:  1 * time.Minute,
		JWTRefreshTokenExpireThrough: 60 * time.Minute,
	}
	user := *testUser()

	whiteListRepo := mocks.NewWhiteListRepository(t)
	whiteListRepo.On(
		"Set",
		ctx,
		user.ID,
		mock.AnythingOfType("string"),
		jwtConfig.JWTRefreshTokenExpireThrough,
	).Return(nil).Once()

	service := jwt.NewService(&logger, jwtConfig, whiteListRepo)
	refreshTokenOld, err := service.GenerateRefreshToken(ctx, user)
	require.NoError(t, err)

	whiteListRepo.On(
		"Get",
		ctx,
		user.ID,
	).Return("new_refresh_token", nil).Once()

	_, err = service.VerifyRefreshToken(ctx, refreshTokenOld)
	require.Error(t, err)
	require.ErrorIs(t, err, autService.ErrMissMatchJWTRefreshWithWhiteList)
}

func TestService_ErrorWhenAccessTokenWithWrongSign(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	jwtConfig := &config.JWTConfig{
		JWTSecretKey:                 "secret",
		JWTAccessTokenExpireThrough:  1 * time.Second,
		JWTRefreshTokenExpireThrough: 60 * time.Minute,
	}
	jwtConfigWithAnotherSign := &config.JWTConfig{
		JWTSecretKey:                 "another_sign",
		JWTAccessTokenExpireThrough:  1 * time.Second,
		JWTRefreshTokenExpireThrough: 60 * time.Minute,
	}
	whiteListRepo := mocks.NewWhiteListRepository(t)

	service := jwt.NewService(&logger, jwtConfig, whiteListRepo)
	serviceWithAnotherSign := jwt.NewService(&logger, jwtConfigWithAnotherSign, whiteListRepo)

	user := *testUser()

	accessTokenWithAnotherSign, err := serviceWithAnotherSign.GenerateAccessToken(user)
	require.NoError(t, err)

	_, err = service.VerifyAccessToken(accessTokenWithAnotherSign)
	require.Error(t, err)
}

func TestService_VerifyAccessTokenFromCtxSuccess(t *testing.T) {
	ctx := context.Background()

	logger := zerolog.Nop()
	jwtConfig := &config.JWTConfig{
		JWTSecretKey:                 "secret",
		JWTAccessTokenExpireThrough:  1 * time.Minute,
		JWTRefreshTokenExpireThrough: 60 * time.Minute,
	}
	whiteListRepo := mocks.NewWhiteListRepository(t)

	service := jwt.NewService(&logger, jwtConfig, whiteListRepo)
	user := *testUser()
	accessToken, err := service.GenerateAccessToken(user)
	require.NoError(t, err)

	md := metadata.New(map[string]string{"Authorization": "Bearer " + accessToken})
	ctx = metadata.NewIncomingContext(ctx, md)

	_, err = service.VerifyAccessTokenFromCtx(ctx)
	require.NoError(t, err)
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
