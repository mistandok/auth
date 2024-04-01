package access

import (
	"context"
	"fmt"

	"github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/internal/service"
)

// Create ..
func (s *Service) Create(ctx context.Context, endpointAccess *model.EndpointAccess) (int64, error) {
	s.logger.Debug().Msg(fmt.Sprintf("попытка создать права для адреса %s с уровнем доступа роли %s", endpointAccess.Address, endpointAccess.Role))
	userClaims, err := s.jwtService.VerifyAccessTokenFromCtx(ctx)
	if err != nil {
		return 0, err
	}

	if userClaims.Role != string(model.ADMIN) {
		return 0, service.ErrNeedAdminRole
	}

	accessID, err := s.endpointAccessRepo.Create(ctx, endpointAccess)
	if err != nil {
		return 0, err
	}

	return accessID, err
}
