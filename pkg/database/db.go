package database

import (
	"database/sql"
	"fmt"

	"github.com/hyperxpizza/auth-service/pkg/config"
	"github.com/hyperxpizza/auth-service/pkg/models"
	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

func Connect(cfg *config.Config) (*Database, error) {
	psqlInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

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

func (db *Database) GetUser(id int64, username string) (*models.User, error) {
	var user models.User
	err := db.QueryRow(`select * from users where id = $1 and username = $2`, id, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Created,
		&user.Updated,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
