package repository

import "errors"

const (
	errMsgUserNotFound           = "пользователь не найден"
	errMsgEmailIsTaken           = "пользователь с таким email уже существует"
	errMsgEndpointAccessExists   = "настройка доступа для адреса с такой ролью уже существует"
	errMsgEndpointAccessNotFound = "настройка доступа для адреса не найден"
	errMsgIncorrectFilters       = "некоректно заданы фильтры"
)

var (
	ErrUserNotFound           = errors.New(errMsgUserNotFound)           // ErrUserNotFound сигнальная ошибка в случае отсутствия пользователя.
	ErrEmailIsTaken           = errors.New(errMsgEmailIsTaken)           // ErrEmailIsTaken сигнальная ошибка в случае дублирования email.
	ErrEndpointAccessExists   = errors.New(errMsgEndpointAccessExists)   // ErrEndpointAccessExists сигнальная ошибка.
	ErrEndpointAccessNotFound = errors.New(errMsgEndpointAccessNotFound) // ErrEndpointAccessNotFound сигнальная ошибка.
	ErrIncorrectFilters       = errors.New(errMsgIncorrectFilters)       // ErrIncorrectFilters сигнальная ошибка.
)
