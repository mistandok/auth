package user

import (
	"context"

	"github.com/mistandok/auth/internal/model"
)

// Delete ..
func (s *Service) Delete(ctx context.Context, userID model.UserID) error {
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		s.logger.Err(err).Msg("не удалось удалить пользователя")
		return err
	}

	return nil
}
