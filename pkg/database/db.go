package database

import (
	"database/sql"

	"github.com/hyperxpizza/auth-service/pkg/config"
)

type Database struct {
	*sql.DB
}

func Connect(cfg *config.Config) (*Database, error) {
	return &Database{}, nil
}
