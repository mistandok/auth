package auth

import (
	"context"
	"errors"
	"github.com/mistandok/auth/internal/service"
	"github.com/mistandok/auth/pkg/auth_v1"
)

const msgInternalError = "что-то пошло не так, мы уже работаем над решением проблемы"

var errInternal = errors.New(msgInternalError)

// Implementation user Server.
type Implementation struct {
	auth_v1.UnimplementedAuthV1Server
	authService service.AuthService
}

// NewImplementation ..
func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}

// Create ..
func (i *Implementation) Login(ctx context.Context, request *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {
	return &auth_v1.LoginResponse{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}, nil
}

// RefreshTokens ..
func (i *Implementation) RefreshTokens(ctx context.Context, request *auth_v1.RefreshTokensRequest) (*auth_v1.RefreshTokensResponse, error) {
	return &auth_v1.RefreshTokensResponse{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}, nil
}
