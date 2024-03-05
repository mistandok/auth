package user

import (
	"context"
	"fmt"

	"github.com/mistandok/auth/internal/model"
)

// Create ..
func (s *Service) Create(ctx context.Context, userForCreate *model.UserForCreate) (int64, error) {
	userID, err := s.userRepo.Create(ctx, userForCreate)
	if err != nil {
		s.logger.Err(err).Msg("не удалось создать пользователя")
		return 0, fmt.Errorf("ошибка при попытке создать пользователя: %w", err)
	}

	return userID, nil
}
