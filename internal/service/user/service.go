package user

import (
	"github.com/mistandok/auth/internal/repository"
	"github.com/mistandok/auth/internal/service"
	"github.com/rs/zerolog"
)

var _ service.UserService = (*Service)(nil)

// Service ..
type Service struct {
	logger   *zerolog.Logger
	userRepo repository.UserRepository
}

// NewService ..
func NewService(logger *zerolog.Logger, userRepository repository.UserRepository) *Service {
	return &Service{
		logger:   logger,
		userRepo: userRepository,
	}
}
