package repositories

import "github.com/pkg/errors"

const (
	errMsgUserNotFound = "пользователь не найден"
	errMsgEmailIsTaken = "пользователь с таким email уже существует"
)

var (
	ErrUserNotFound = errors.New(errMsgUserNotFound) // ErrUserNotFound сигнальная ошибка в случае отсутствия пользователя.
	ErrEmailIsTaken = errors.New(errMsgEmailIsTaken) // ErrEmailIsTaken сигнальная ошибка в случае дублирования email.
)
