package database

import (
	"database/sql"
	"fmt"
	"time"

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

func (db *Database) InsertUser(user models.User) (int64, error) {
	var id int64
	stmt, err := db.Prepare(`insert into users (id, username, email, passwordHash, created, updated) values (default, $1, $2, $3, $4) returning id`)
	if err != nil {
		return 0, err
	}

	err = stmt.QueryRow(user.Username, user.PasswordHash, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Database) DeleteUser(id int64) error {
	stmt, err := db.Prepare(`delete from users where id = $1`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateUser(user models.User) error {
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
