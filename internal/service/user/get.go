package user

import (
	"context"
	"github.com/mistandok/auth/internal/utils"

	"github.com/mistandok/auth/internal/model"
)

// Get ..
func (s *Service) Get(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.userRepo.GetByFilter(ctx, &model.UserFilter{ID: utils.Pointer[int64](userID)})
	if err != nil {
		s.logger.Error().Err(err).Msg("не удалось получить пользователя")
		return nil, err
	}

	return user, nil
}
