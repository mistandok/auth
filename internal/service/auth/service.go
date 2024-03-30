package auth

import (
	"github.com/mistandok/auth/internal/repository"
	"github.com/mistandok/auth/internal/service"
	"github.com/rs/zerolog"
)

var _ service.AuthService = (*Service)(nil)

// Service ..
type Service struct {
	logger   *zerolog.Logger
	userRepo repository.UserRepository
}

// NewService ..
func NewService(logger *zerolog.Logger, userRepo repository.UserRepository) *Service {
	return &Service{
		logger:   logger,
		userRepo: userRepo,
	}
}
