package jwt

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/dgrijalva/jwt-go"
	"github.com/mistandok/auth/internal/config"
	"github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/internal/repository"
	"github.com/mistandok/auth/internal/service"
	"github.com/rs/zerolog"
)

// Service сервис для работы с токенами.
type Service struct {
	logger        *zerolog.Logger
	jwtConfig     *config.JWTConfig
	whiteListRepo repository.WhiteListRepository
}

// NewService новый сервис для работы с токенами.
func NewService(logger *zerolog.Logger, jwtConfig *config.JWTConfig, whiteListRepo repository.WhiteListRepository) *Service {
	return &Service{
		logger:        logger,
		jwtConfig:     jwtConfig,
		whiteListRepo: whiteListRepo,
	}
}

// GenerateAccessToken генерация access токена.
func (s *Service) GenerateAccessToken(user model.User) (string, error) {
	s.logger.Debug().Int64("userID", user.ID).Msg("попытка сгенерировать access token")
	return s.generateToken(user, s.jwtConfig.JWTAccessTokenExpireThrough)
}

// GenerateRefreshToken генерация refresh токена.
func (s *Service) GenerateRefreshToken(ctx context.Context, user model.User) (string, error) {
	s.logger.Debug().Int64("userID", user.ID).Msg("попытка сгенерировать refresh token")

	token, err := s.generateToken(user, s.jwtConfig.JWTRefreshTokenExpireThrough)
	if err != nil {
		return "", fmt.Errorf("ошибка при генерации refresh token: %w", err)
	}

	err = s.whiteListRepo.Set(ctx, user.ID, token, s.jwtConfig.JWTRefreshTokenExpireThrough)
	if err != nil {
		s.logger.Error().Err(err).Msg("не удалось сохранить токен в белом списке")
		return "", err
	}

	return token, nil
}

// VerifyAccessToken проверка access токена на валидность.
func (s *Service) VerifyAccessToken(tokenStr string) (*model.UserClaims, error) {
	s.logger.Debug().Str("access_token", tokenStr).Msg("верификация access token")
	return s.verifyToken(tokenStr)
}

// VerifyAccessTokenFromCtx проверка access токена из контекста.
func (s *Service) VerifyAccessTokenFromCtx(ctx context.Context) (*model.UserClaims, error) {
	accessToken, err := s.fetchTokenFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return s.verifyToken(accessToken)
}

// VerifyRefreshToken проверка refresh токена на валидность.
func (s *Service) VerifyRefreshToken(ctx context.Context, tokenStr string) (*model.UserClaims, error) {
	s.logger.Debug().Str("refresh_token", tokenStr).Msg("верификация refresh token")
	userClaims, err := s.verifyToken(tokenStr)
	if err != nil {
		return nil, err
	}

	userID, err := strconv.ParseInt(userClaims.UserID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить ID пользователя: %w", err)
	}

	actualToken, err := s.whiteListRepo.Get(ctx, userID)
	if err != nil {
		s.logger.Error().Err(nil).Msg("не удалось получить токен из бeлого списка")
		return nil, fmt.Errorf("не удалось получить токен для пользователя %d из белого списка: %w", userID, err)
	}

	if actualToken != tokenStr {
		return nil, service.ErrMissMatchJWTRefreshWithWhiteList
	}

	return userClaims, nil
}

func (s *Service) generateToken(user model.User, duration time.Duration) (string, error) {
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		UserID:   fmt.Sprintf("%d", user.ID),
		UserName: user.Name,
		Role:     string(user.Role),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtConfig.JWTSecretKey))
	if err != nil {
		return "", fmt.Errorf("не удалось подписать токен: %w", err)
	}

	return signedToken, nil
}

func (s *Service) verifyToken(tokenStr string) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("неожиданный метод подписи токена")
			}

			return []byte(s.jwtConfig.JWTSecretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("невалидный токен: %w", err)
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, fmt.Errorf("невалидное тело токена")
	}

	return claims, nil
}

func (s *Service) fetchTokenFromCtx(ctx context.Context) (string, error) {
	s.logger.Debug().Msg("попытка достать access token из ctx")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("метаданные не переданы")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return "", errors.New("в header не представлена authorization")
	}

	if !strings.HasPrefix(authHeader[0], "Bearer ") {
		return "", errors.New("некоректный формат authorization в header")
	}

	return strings.TrimPrefix(authHeader[0], "Bearer "), nil
}
