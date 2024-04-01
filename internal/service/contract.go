package service

import (
	"context"

	"github.com/mistandok/auth/internal/model"
)

//go:generate ../../bin/mockery --output ./mocks  --inpackage-suffix --all --case snake

// UserService ..
type UserService interface {
	Create(context.Context, *model.UserForCreate) (int64, error)
	Update(context.Context, *model.UserForUpdate) error
	Get(context.Context, int64) (*model.User, error)
	Delete(context.Context, int64) error
}

// AuthService ..
type AuthService interface {
	Login(ctx context.Context, userEmail string, userPassword string) (*model.Tokens, error)
	RefreshTokens(ctx context.Context, refreshTokenStr string) (*model.Tokens, error)
}

// JWTService ..
type JWTService interface {
	GenerateAccessToken(user model.User) (string, error)
	GenerateRefreshToken(ctx context.Context, user model.User) (string, error)
	VerifyAccessToken(tokenStr string) (*model.UserClaims, error)
	VerifyAccessTokenFromCtx(ctx context.Context) (*model.UserClaims, error)
	VerifyRefreshToken(ctx context.Context, tokenStr string) (*model.UserClaims, error)
}

// AccessService ..
type AccessService interface {
	Create(context.Context, *model.EndpointAccess) (int64, error)
	Check(ctx context.Context, endpointAddress string) error
}
