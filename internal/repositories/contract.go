package repositories

import (
	"context"
	"time"
)

// CRUDUserCreateIn params for create.
type CRUDUserCreateIn struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// CRUDUserCreateOut out params for create.
type CRUDUserCreateOut struct {
	ID int64 `json:"id"`
}

// CRUDUserUpdateIn params for uodate
type CRUDUserUpdateIn struct {
	ID    int64   `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Role  *string `json:"role"`
}

// CRUDUserGetIn params for get.
type CRUDUserGetIn struct {
	ID int64 `json:"id"`
}

// CRUDUserGetOut out params for get.
type CRUDUserGetOut struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CRUDUserDeleteIn params for delete
type CRUDUserDeleteIn struct {
	ID int64 `json:"id"`
}

// CRUDUserRepository interface for crud user repositories
type CRUDUserRepository interface {
	Create(context.Context, *CRUDUserCreateIn) (*CRUDUserCreateOut, error)
	Update(context.Context, *CRUDUserUpdateIn) error
	Get(context.Context, *CRUDUserGetIn) (*CRUDUserGetOut, error)
	Delete(context.Context, *CRUDUserDeleteIn) error
}
