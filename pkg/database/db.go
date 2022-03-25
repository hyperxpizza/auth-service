package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/hyperxpizza/auth-service/pkg/config"
	pb "github.com/hyperxpizza/auth-service/pkg/grpc"
	_ "github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (db *Database) InsertUser(user *pb.AuthServiceUser) (int64, error) {
	var id int64
	stmt, err := db.Prepare(`insert into users (id, username, passwordHash, created, updated, relatedUsersServiceID) values (default, $1, $2, $3, $4, $5) returning id`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(user.Username, user.PasswordHash, time.Now(), time.Now(), user.RelatedUsersServiceID).Scan(&id)
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
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateUser(user *pb.AuthServiceUser) error {
	stmt, err := db.Prepare(`update users set username=$1, passwordHash=$2, updated=$3 where id=$4`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.PasswordHash, time.Now(), user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetUserByUsersServiceID(id int64, username string) (*pb.AuthServiceUser, error) {
	var user pb.AuthServiceUser
	var created time.Time
	var updated time.Time

	err := db.QueryRow(`select * from users where relatedUsersServiceID = $1 and username = $2`, id, username).Scan(
		&user.Id,
		&user.Username,
		&user.PasswordHash,
		&created,
		&updated,
		&user.RelatedUsersServiceID,
	)
	if err != nil {
		return nil, err
	}

	user.Created = timestamppb.New(created)
	user.Updated = timestamppb.New(updated)

	return &user, nil
}

func (db *Database) GetUser(id int64, username string) (*pb.AuthServiceUser, error) {
	var user pb.AuthServiceUser
	var created time.Time
	var updated time.Time

	err := db.QueryRow(`select * from users where id = $1 and username = $2`, id, username).Scan(
		&user.Id,
		&user.Username,
		&user.PasswordHash,
		&created,
		&updated,
		&user.RelatedUsersServiceID,
	)
	if err != nil {
		return nil, err
	}

	user.Created = timestamppb.New(created)
	user.Updated = timestamppb.New(updated)

	return &user, nil
}

func (db *Database) ChangePassword(id int64, username, passwordHash string) error {
	stmt, err := db.Prepare(`update users set passwordHash=$1, updated=$2 where id=$3 and username=$4`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(passwordHash, time.Now(), id, username)
	if err != nil {
		return err
	}

	return nil
}
