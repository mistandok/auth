package model

import "time"

type UserID int64

type UserName string

type UserEmail string

type UserRole string

type Password string

const (
	UNKNOWN UserRole = "UNKNOWN"
	ADMIN   UserRole = "ADMIN"
	USER    UserRole = "USER"
)

type User struct {
	ID        UserID
	Name      UserName
	Email     UserEmail
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserForCreate struct {
	Name     UserName
	Email    UserEmail
	Password Password
	Role     UserRole
}

type UserForUpdate struct {
	ID    UserID
	Name  *UserName
	Email *UserEmail
	Role  *UserRole
}
