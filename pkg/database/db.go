package database

import (
	"database/sql"
	"fmt"

	"github.com/hyperxpizza/auth-service/pkg/config"
	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

func Connect(cfg *config.Config) (*Database, error) {
	psqlInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name)

	database, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{database}, nil
}

func (db *Database) InsertUser() error {
	return nil
}

func (db *Database) DelteUser() error {
	return nil
}

func (db *Database) UpdateUser() error {
	return nil
}
