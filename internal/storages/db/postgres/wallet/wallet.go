package wallet

import (
	"context"
	"fmt"
	"itk-wallet/internal/model"
	"itk-wallet/internal/storages/db/postgres"
)

type Wallet struct {
	db postgres.DB
}

func NewWallet(db postgres.DB) *Wallet {
	return &Wallet{db: db}
}

func (r *Wallet) Insert(ctx context.Context, wallet model.Wallet) error {
	query := `INSERT INTO wallets (balance, uuid)
				VALUES ($1, $2)`
	if _, err := r.db.Exec(ctx, query, wallet.Balance, wallet.UUID); err != nil {
		return fmt.Errorf("insert wallet: %w", err)
	}
	return nil
}

func (r *Wallet) Update(ctx context.Context, wallet model.Wallet) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	selectQuery := `SELECT balance FROM wallets WHERE uuid = $1 FOR UPDATE`
	var currentBalance int64
	if err = tx.QueryRow(ctx, selectQuery, wallet.UUID).Scan(&currentBalance); err != nil {
		return fmt.Errorf("select balance: %w", err)
	}

	updateQuery := `
					UPDATE  wallets
					SET balance = $1
					WHERE uuid = $2
		`
	if _, err = tx.Exec(ctx, updateQuery, wallet.Balance, wallet.UUID); err != nil {
		return fmt.Errorf("update wallet: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

func (r *Wallet) Get(ctx context.Context, uuid int64) (int64, error) {
	query := `SELECT balance
			FROM wallets
			WHERE uuid = $1`

	var amount int64
	if err := r.db.QueryRow(ctx, query, uuid).Scan(&amount); err != nil {
		return 0, fmt.Errorf("get wallet: %w", err)
	}
	return amount, nil
}
