package user

import (
	"context"
	"fmt"
	"github.com/mistandok/auth/internal/model"
)

func (s *Service) Create(ctx context.Context, userForCreate *model.UserForCreate) (model.UserID, error) {
	userID, err := s.userRepo.Create(ctx, userForCreate)
	if err != nil {
		s.logger.Err(err).Msg("не удалось создать пользователя")
		return 0, fmt.Errorf("ошибка при попытке создать пользователя: %w", err)
	}

	return userID, nil
}
