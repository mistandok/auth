package model

import "time"

// UserID идентификатор пользователя
type UserID int64

// UserName имя пользователя
type UserName string

// UserEmail email пользователя
type UserEmail string

// UserRole роль пользователя
type UserRole string

// Password пароль пользователя
type Password string

const (
	UNKNOWN UserRole = "UNKNOWN" // UNKNOWN неизвестный тип пользователя
	ADMIN   UserRole = "ADMIN"   // ADMIN админ
	USER    UserRole = "USER"    // USER обычный пользователь
)

// User ..
type User struct {
	ID        UserID
	Name      UserName
	Email     UserEmail
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserForCreate ..
type UserForCreate struct {
	Name     UserName
	Email    UserEmail
	Password Password
	Role     UserRole
}

// UserForUpdate ..
type UserForUpdate struct {
	ID    UserID
	Name  *UserName
	Email *UserEmail
	Role  *UserRole
}
