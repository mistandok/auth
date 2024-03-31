package auth

import (
	"context"
	"errors"

	"github.com/mistandok/auth/internal/repository"
	"github.com/mistandok/auth/internal/service"
	"github.com/mistandok/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// Login ..
func (i *Implementation) Login(ctx context.Context, request *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {
	tokens, err := i.authService.Login(ctx, request.Email, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, repository.ErrUserNotFound.Error())
		case errors.Is(err, service.ErrIncorrectPassword):
			return nil, status.Error(codes.InvalidArgument, service.ErrIncorrectPassword.Error())
		default:
			return nil, errInternal
		}
	}

	return &auth_v1.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

// RefreshTokens ..
func (i *Implementation) RefreshTokens(ctx context.Context, request *auth_v1.RefreshTokensRequest) (*auth_v1.RefreshTokensResponse, error) {
	tokens, err := i.authService.RefreshTokens(ctx, request.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &auth_v1.RefreshTokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
