package model

import "time"

type User struct {
	ID        int64     `db:"id"`
	Username  string    `db:"user_name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
