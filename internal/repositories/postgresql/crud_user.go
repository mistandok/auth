package postgresql

import (
	"github.com/jackc/pgx/v5"
	"github.com/mistandok/auth/internal/repositories"
	"github.com/rs/zerolog"
)

type CRUDUser struct {
	connect *pgx.Conn
	logger  *zerolog.Logger
}

func NewCRUDUser(conn *pgx.Conn, logger *zerolog.Logger) *CRUDUser {
	return &CRUDUser{
		connect: conn,
		logger:  logger,
	}
}

func (u *CRUDUser) Create(in *repositories.CRUDUserCreateIn) (*repositories.CRUDUserCreateOut, error) {
	return nil, nil
}

func (u *CRUDUser) Update(in *repositories.CRUDUserUpdateIn) error {
	return nil
}

func (u *CRUDUser) Get(in *repositories.CRUDUserGetIn) (*repositories.CRUDUserGetOut, error) {
	return nil, nil
}

func (u *CRUDUser) Delete(in *repositories.CRUDUserDeleteIn) error {
	return nil
}
