package repositories

import (
	"context"
	"time"
)

type CRUDUserCreateIn struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type CRUDUserCreateOut struct {
	Id int64 `json:"id"`
}

type CRUDUserUpdateIn struct {
	Id    int64   `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Role  *string `json:"role"`
}

type CRUDUserGetIn struct {
	Id int64 `json:"id"`
}

type CRUDUserGetOut struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CRUDUserDeleteIn struct {
	Id int64 `json:"id"`
}

type CRUDUserRepository interface {
	Create(context.Context, *CRUDUserCreateIn) (*CRUDUserCreateOut, error)
	Update(context.Context, *CRUDUserUpdateIn) error
	Get(context.Context, *CRUDUserGetIn) (*CRUDUserGetOut, error)
	Delete(context.Context, *CRUDUserDeleteIn) error
}
