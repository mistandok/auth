package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// UserEmail email пользователя
type UserEmail string

// UserRole роль пользователя
type UserRole string

const (
	UNKNOWN UserRole = "UNKNOWN" // UNKNOWN неизвестный тип пользователя
	ADMIN   UserRole = "ADMIN"   // ADMIN админ
	USER    UserRole = "USER"    // USER обычный пользователь
)

// User ..
type User struct {
	ID        int64
	Name      string
	Email     UserEmail
	Role      UserRole
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserForCreate ..
type UserForCreate struct {
	Name     string
	Email    UserEmail
	Password string
	Role     UserRole
}

// UserForUpdate ..
type UserForUpdate struct {
	ID    int64
	Name  *string
	Email *UserEmail
	Role  *UserRole
}

// UserClaims ..
type UserClaims struct {
	jwt.StandardClaims
	UserID   string `json:"userID"`
	UserName string `json:"username"`
	Role     string `json:"role"`
}
