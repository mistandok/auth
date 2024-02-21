package repositories

import "time"

type CRUDUserCreateIn struct {
	Name     string
	Email    string
	Password string
	Role     string
}

type CRUDUserCreateOut struct {
	Id int64
}

type CRUDUserUpdateIn struct {
	Id    int64
	Name  string
	Email string
	Role  string
}

type CRUDUserGetIn struct {
	Id int64
}

type CRUDUserGetOut struct {
	Id        int64
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CRUDUserDeleteIn struct {
	Id int64
}

type CRUDUserExecutor interface {
	Create(*CRUDUserCreateIn) (*CRUDUserCreateOut, error)
	Update(*CRUDUserUpdateIn) error
	Get(*CRUDUserGetIn) (*CRUDUserGetOut, error)
	Delete(*CRUDUserDeleteIn) error
}
