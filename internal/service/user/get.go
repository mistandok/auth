package user

import (
	"context"

	"github.com/mistandok/auth/internal/model"
)

// Get ..
func (s *Service) Get(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		s.logger.Err(err).Msg("не удалось получить пользователя")
		return nil, err
	}

	return user, nil
}
