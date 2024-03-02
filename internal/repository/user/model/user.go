package model

import "time"

type User struct {
	ID        int64
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserForCreate struct {
	Name     string
	Email    string
	Password string
	Role     string
}

type UserForUpdate struct {
	ID    int64
	Name  *string
	Email *string
	Role  *string
}
