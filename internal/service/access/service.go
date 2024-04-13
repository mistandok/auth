package access

import (
	"github.com/mistandok/auth/internal/repository"
	"github.com/mistandok/auth/internal/service"
	"github.com/rs/zerolog"
)

var _ service.AccessService = (*Service)(nil)

// Service ..
type Service struct {
	logger             *zerolog.Logger
	endpointAccessRepo repository.EndpointAccessRepository
	jwtService         service.JWTService
}

// NewService ..
func NewService(logger *zerolog.Logger, endpointAccessRepo repository.EndpointAccessRepository, jwtService service.JWTService) *Service {
	return &Service{
		logger:             logger,
		endpointAccessRepo: endpointAccessRepo,
		jwtService:         jwtService,
	}
}
