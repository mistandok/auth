package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mistandok/auth/internal/config"
	"github.com/mistandok/auth/internal/model"
	"time"
)

type Manager struct {
	jwtConfig *config.JWTConfig
}

func NewManager(jwtConfig *config.JWTConfig) *Manager {
	return &Manager{jwtConfig: jwtConfig}
}

func (m *Manager) VerifyToken(tokenStr string) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("неожиданный метод подписи токена")
			}

			return []byte(m.jwtConfig.JWTSecretKey), nil
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

func (m *Manager) GenerateAccessToken(user model.User) (string, error) {
	return m.generateToken(user, m.jwtConfig.JWTAccessTokenExpireThrough)
}

func (m *Manager) GenerateRefreshToken(user model.User) (string, error) {
	return m.generateToken(user, m.jwtConfig.JWTRefreshTokenExpireThrough)
}

func (m *Manager) generateToken(user model.User, duration time.Duration) (string, error) {
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		UserID:   fmt.Sprintf("%d", user.ID),
		UserName: user.Name,
		Role:     string(user.Role),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(m.jwtConfig.JWTSecretKey))
	if err != nil {
		return "", fmt.Errorf("не удалось подписать токен: %w", err)
	}

	return signedToken, nil
}
