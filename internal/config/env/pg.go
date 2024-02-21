package env

import (
	"github.com/mistandok/auth/internal/config"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

const (
	host     = "PG_HOST"
	port     = "PG_PORT"
	user     = "POSTGRES_USER"
	password = "POSTGRES_PASSWORD"
	dbName   = "POSTGRES_DB"
)

type PgCfgSearcher struct{}

func NewPgCfgSearcher() *PgCfgSearcher {
	return &PgCfgSearcher{}
}

func (s *PgCfgSearcher) Get() (*config.PgConfig, error) {
	dbHost := os.Getenv(host)
	if len(dbHost) == 0 {
		return nil, errors.New("db host not found")
	}

	dbPort := os.Getenv(port)
	if len(dbPort) == 0 {
		return nil, errors.New("db port not found")
	}

	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		return nil, errors.Errorf("некорректный формат порта: %v", err)
	}

	dbUser := os.Getenv(user)
	if len(dbUser) == 0 {
		return nil, errors.New("db user not found")
	}

	dbPass := os.Getenv(password)
	if len(dbPass) == 0 {
		return nil, errors.New("db password not found")
	}

	name := os.Getenv(dbName)
	if len(name) == 0 {
		return nil, errors.New("db name not found")
	}

	return &config.PgConfig{
		Host:     dbHost,
		Port:     dbPortInt,
		User:     dbUser,
		Password: dbPass,
		DbName:   name,
	}, nil
}
