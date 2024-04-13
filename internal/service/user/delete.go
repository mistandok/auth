package user

import (
	"context"
)

// Delete ..
func (s *Service) Delete(ctx context.Context, userID int64) error {
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		s.logger.Error().Err(err).Msg("не удалось удалить пользователя")
		return err
	}

	return nil
}
