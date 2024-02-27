package repositories

import (
	"context"
	"time"
)

// UserRepository interface for crud user repositories
type UserRepository interface {
	Create(context.Context, *UserCreateIn) (*UserCreateOut, error)
	Update(context.Context, *UserUpdateIn) error
	Get(context.Context, *UserGetIn) (*UserGetOut, error)
	Delete(context.Context, *UserDeleteIn) error
}

// UserCreateIn params for create.
type UserCreateIn struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UserCreateOut out params for create.
type UserCreateOut struct {
	ID int64 `json:"id"`
}

// UserUpdateIn params for uodate
type UserUpdateIn struct {
	ID    int64   `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Role  *string `json:"role"`
}

// UserGetIn params for get.
type UserGetIn struct {
	ID int64 `json:"id"`
}

// UserGetOut out params for get.
type UserGetOut struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UserDeleteIn params for delete
type UserDeleteIn struct {
	ID int64 `json:"id"`
}
