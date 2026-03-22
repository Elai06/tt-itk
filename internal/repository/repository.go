package repository

import (
	"context"
	"fmt"
	"itk/internal/dto"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

//go:generate mockgen -destination=mocks/mock_repository.go -package=mocks itk/internal/repository Repository

type Repository interface {
	Create(ctx context.Context, wallet dto.WalletRequest) error
	Update(ctx context.Context, wallet dto.WalletRequest) error
	Get(ctx context.Context, uuid int64) (int64, error)
	GetConConfig() *pgx.ConnConfig
	Close(ctx context.Context) error
}

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close(ctx context.Context) error
	Config() *pgx.ConnConfig
	Ping(ctx context.Context) error
}

type repository struct {
	Con DB
}

func NewRepository(dbUrl string, ctx context.Context) (Repository, error) {
	con, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = con.Ping(ctx)
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	return &repository{Con: con}, err
}

func (s *repository) Create(ctx context.Context, wallet dto.WalletRequest) error {
	query := `INSERT INTO wallets (balance, uuid) 
				VALUES ($1, $2)`
	res, err := s.Con.Exec(ctx, query, wallet.Amount, wallet.UUID)
	if err != nil {
		log.Println("Error inserting wallet", err)
		return err
	}

	log.Println("inserted", res.RowsAffected())

	return nil
}

func (s *repository) Update(ctx context.Context, wallet dto.WalletRequest) error {
	query := `
				UPDATE  wallets 
				SET balance = $1
				WHERE uuid = $2
	`
	res, err := s.Con.Exec(ctx, query, wallet.Amount, wallet.UUID)
	if err != nil {
		log.Println("Error inserting wallet", err)
		return err
	}

	log.Println("inserted", res.RowsAffected())

	return nil
}

func (s *repository) Get(ctx context.Context, uuid int64) (int64, error) {
	query := `SELECT balance
			FROM wallets 
			WHERE uuid = $1`

	var amount int64
	rows := s.Con.QueryRow(ctx, query, uuid)
	err := rows.Scan(&amount)
	if err != nil {
		return 0, err
	}

	log.Println("balance", amount)

	return amount, nil
}

func (s *repository) GetConConfig() *pgx.ConnConfig {
	return s.Con.Config().Copy()
}

func (s *repository) Close(ctx context.Context) error {
	return s.Con.Close(ctx)
}
