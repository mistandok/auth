package access

import (
	"context"

	"github.com/mistandok/auth/internal/model"
)

// Check ..
func (s *Service) Check(ctx context.Context, endpointAddress string) error {
	userClaims, err := s.jwtService.VerifyAccessTokenFromCtx(ctx)
	if err != nil {
		return err
	}

	if userClaims.Role == string(model.ADMIN) {
		return nil
	}

	_, err = s.endpointAccessRepo.GetByAddressAndRole(
		ctx,
		endpointAddress,
		model.UserRole(userClaims.Role),
	)
	if err != nil {
		return err
	}

	return nil
}
