package user

import (
	"context"
	"fmt"
	"github.com/mistandok/auth/internal/model"
)

func (s *Service) Update(ctx context.Context, userForUpdate *model.UserForUpdate) error {
	if err := s.userRepo.Update(ctx, userForUpdate); err != nil {
		s.logger.Err(err).Msg("не удалось обновить пользователя")
		return fmt.Errorf("ошибка при попытке обновить пользователя: %w", err)
	}

	return nil
}
