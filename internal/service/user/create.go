package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/mistandok/auth/internal/model"
)

var ErrPassToLong = errors.New("слишком длинный пароль") // ErrPassToLong ..

// Create ..
func (s *Service) Create(ctx context.Context, userForCreate *model.UserForCreate) (int64, error) {
	hashedPassword, err := s.passManager.HashPassword(userForCreate.Password)
	if err != nil {
		s.logger.Error().Err(err).Msg("не удалось хэшировать пароль")
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return 0, ErrPassToLong
		}
		return 0, err
	}

	userForCreate.Password = hashedPassword

	userID, err := s.userRepo.Create(ctx, userForCreate)
	if err != nil {
		s.logger.Error().Err(err).Msg("не удалось создать пользователя")
		return 0, fmt.Errorf("ошибка при попытке создать пользователя: %w", err)
	}

	return userID, nil
}
