package service

import "errors"

const (
	errMsgIncorrectPass                    = "неправильный пароль" //nolint:gosec,unused
	errMsgMissMatchJWTRefreshWithWhiteList = "refresh token не совпадает с refresh token в белом списке"
)

var (
	ErrIncorrectPassword                = errors.New(errMsgIncorrectPass)                    // ErrIncorrectPassword сигнальная ошибка в случае несовпадения паролей.
	ErrMissMatchJWTRefreshWithWhiteList = errors.New(errMsgMissMatchJWTRefreshWithWhiteList) // ErrMissMatchJWTRefreshWithWhiteList сигнальная ошибка о том, что refresh токен не совпадает с refresh токеном из белого списка.
)
