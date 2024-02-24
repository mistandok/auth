package repositories

import "github.com/pkg/errors"

const (
	errMsgUserNotFound = "пользователь не найден"
	errMsgEmailIsTaken = "пользователь с таким email уже существует"
)

var (
	ErrUserNotFound = errors.New(errMsgUserNotFound)
	ErrEmailIsTaken = errors.New(errMsgEmailIsTaken)
)
