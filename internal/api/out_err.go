package api

import "errors"

const msgInternalError = "что-то пошло не так, мы уже работаем над решением проблемы"

var ErrInternal = errors.New(msgInternalError) // ErrInternal внутренняя ошибка.
