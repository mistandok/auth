package auth

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mistandok/auth/internal/utils"

	"github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/internal/repository"
	"github.com/mistandok/auth/internal/service"
	"github.com/mistandok/auth/internal/utils/password"
	"github.com/rs/zerolog"
)

var _ service.AuthService = (*Service)(nil)

// Service ..
type Service struct {
	logger      *zerolog.Logger
	userRepo    repository.UserRepository
	jwtService  service.JWTService
	passManager *password.Manager
}

// NewService ..
func NewService(
	logger *zerolog.Logger,
	userRepo repository.UserRepository,
	jwtService service.JWTService,
	passManager *password.Manager,
) *Service {
	return &Service{
		logger:      logger,
		userRepo:    userRepo,
		jwtService:  jwtService,
		passManager: passManager,
	}
}

// Login пользователя по email и паролю.
func (s *Service) Login(ctx context.Context, userEmail string, userPassword string) (*model.Tokens, error) {
	user, err := s.userRepo.GetByFilter(ctx, &model.UserFilter{Email: utils.Pointer[string](userEmail)})
	if err != nil {
		return nil, err
	}

	if !s.passManager.CheckPasswordHash(userPassword, user.Password) {
		return nil, service.ErrIncorrectPassword
	}

	accessToken, err := s.jwtService.GenerateAccessToken(*user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(ctx, *user)
	if err != nil {
		return nil, err
	}

	return &model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshTokens обновление пары токенов по refresh token
func (s *Service) RefreshTokens(ctx context.Context, refreshTokenStr string) (*model.Tokens, error) {
	userClaims, err := s.jwtService.VerifyRefreshToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, err
	}

	user, err := userFromUserClaims(userClaims)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtService.GenerateAccessToken(*user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(ctx, *user)
	if err != nil {
		return nil, err
	}

	return &model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func userFromUserClaims(claims *model.UserClaims) (*model.User, error) {
	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить ID пользователя ин UserClaims: %w", err)
	}

	return &model.User{
		ID:   userID,
		Name: claims.UserName,
		Role: model.UserRole(claims.Role),
	}, nil
}
