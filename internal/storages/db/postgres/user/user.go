package user

import (
	"context"
	"fmt"
	"itk-wallet/internal/model"
	"itk-wallet/internal/storages/db/postgres"
)

type User struct {
	db postgres.DB
}

func NewUser(db postgres.DB) *User {
	return &User{db: db}
}

func (u *User) Insert(ctx context.Context, user model.User) error {
	query := `INSERT INTO users (user_name, email, password)
				VALUES ($1, $2, $3)`

	if _, err := u.db.Exec(ctx, query, user.Username, user.Email, user.Password); err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (u *User) Get(ctx context.Context, email string) (model.User, error) {
	var usrModel model.User
	query := `SELECT id, user_name, email, password, created_at, updated_at
				FROM users
				WHERE email = $1`
	err := u.db.QueryRow(ctx, query, email).Scan(
		&usrModel.ID,
		&usrModel.Username,
		&usrModel.Email,
		&usrModel.Password,
		&usrModel.CreatedAt,
		&usrModel.UpdatedAt,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("get user: %w", err)
	}
	return usrModel, nil
}
