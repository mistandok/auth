package user

import (
	"github.com/mistandok/auth/internal/repository"
	"github.com/mistandok/auth/internal/service"
	"github.com/mistandok/auth/internal/utils/password"
	"github.com/rs/zerolog"
)

var _ service.UserService = (*Service)(nil)

// Service ..
type Service struct {
	logger      *zerolog.Logger
	userRepo    repository.UserRepository
	passManager *password.Manager
}

// NewService ..
func NewService(logger *zerolog.Logger, userRepository repository.UserRepository, passManager *password.Manager) *Service {
	return &Service{
		logger:      logger,
		userRepo:    userRepository,
		passManager: passManager,
	}
}
