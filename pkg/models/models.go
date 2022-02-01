package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"passwordHash"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}
