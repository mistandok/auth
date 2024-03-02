package model

import "time"

// User ..
type User struct {
	ID        int64
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserForCreate ..
type UserForCreate struct {
	Name     string
	Email    string
	Password string
	Role     string
}

// UserForUpdate ..
type UserForUpdate struct {
	ID    int64
	Name  *string
	Email *string
	Role  *string
}
