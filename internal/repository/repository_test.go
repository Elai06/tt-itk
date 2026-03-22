package repository

import (
	"context"
	"errors"
	"itk/internal/dto"
	"regexp"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v5"
)

func TestRepository_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		wallet    dto.WalletRequest
		mockSetup func(mock pgxmock.PgxConnIface)
		wantErr   bool
	}{
		{
			name: "success",
			wallet: dto.WalletRequest{
				UUID:   1,
				Amount: 100,
			},
			mockSetup: func(mock pgxmock.PgxConnIface) {
				query := regexp.QuoteMeta(`INSERT INTO wallets (balance, uuid) 
				VALUES ($1, $2)`)
				mock.ExpectExec(query).
					WithArgs(int64(100), int64(1)).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			wantErr: false,
		},
		{
			name: "db error",
			wallet: dto.WalletRequest{
				UUID:   2,
				Amount: 200,
			},
			mockSetup: func(mock pgxmock.PgxConnIface) {
				query := regexp.QuoteMeta(`INSERT INTO wallets (balance, uuid) 
				VALUES ($1, $2)`)
				mock.ExpectExec(query).
					WithArgs(int64(200), int64(2)).
					WillReturnError(errors.New("insert failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewConn()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}

			tt.mockSetup(mock)

			repo := &repository{Con: mock}

			err = repo.Create(ctx, tt.wallet)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

func TestRepository_Update(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		wallet    dto.WalletRequest
		mockSetup func(mock pgxmock.PgxConnIface)
		wantErr   bool
	}{
		{
			name: "success",
			wallet: dto.WalletRequest{
				UUID:   1,
				Amount: 500,
			},
			mockSetup: func(mock pgxmock.PgxConnIface) {
				query := regexp.QuoteMeta(`
				UPDATE  wallets 
				SET balance = $1
				WHERE uuid = $2
	`)
				mock.ExpectExec(query).
					WithArgs(int64(500), int64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: false,
		},
		{
			name: "db error",
			wallet: dto.WalletRequest{
				UUID:   2,
				Amount: 700,
			},
			mockSetup: func(mock pgxmock.PgxConnIface) {
				query := regexp.QuoteMeta(`
				UPDATE  wallets 
				SET balance = $1
				WHERE uuid = $2
	`)
				mock.ExpectExec(query).
					WithArgs(int64(700), int64(2)).
					WillReturnError(errors.New("update failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewConn()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}

			tt.mockSetup(mock)

			repo := &repository{Con: mock}

			err = repo.Update(ctx, tt.wallet)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

func TestRepository_Get(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name       string
		uuid       int64
		mockSetup  func(mock pgxmock.PgxConnIface)
		wantAmount int64
		wantErr    bool
	}{
		{
			name: "success",
			uuid: 1,
			mockSetup: func(mock pgxmock.PgxConnIface) {
				rows := pgxmock.NewRows([]string{"balance"}).AddRow(int64(1000))
				query := regexp.QuoteMeta(`SELECT balance
			FROM wallets 
			WHERE uuid = $1`)
				mock.ExpectQuery(query).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			wantAmount: 1000,
			wantErr:    false,
		},
		{
			name: "not found",
			uuid: 2,
			mockSetup: func(mock pgxmock.PgxConnIface) {
				query := regexp.QuoteMeta(`SELECT balance
			FROM wallets 
			WHERE uuid = $1`)
				mock.ExpectQuery(query).
					WithArgs(int64(2)).
					WillReturnError(pgx.ErrNoRows)
			},
			wantAmount: 0,
			wantErr:    true,
		},
		{
			name: "scan error",
			uuid: 3,
			mockSetup: func(mock pgxmock.PgxConnIface) {
				rows := pgxmock.NewRows([]string{"balance"}).AddRow("bad_value")
				query := regexp.QuoteMeta(`SELECT balance
			FROM wallets 
			WHERE uuid = $1`)
				mock.ExpectQuery(query).
					WithArgs(int64(3)).
					WillReturnRows(rows)
			},
			wantAmount: 0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewConn()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}

			tt.mockSetup(mock)

			repo := &repository{Con: mock}

			got, err := repo.Get(ctx, tt.uuid)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.wantAmount {
				t.Fatalf("Get() got = %d, want %d", got, tt.wantAmount)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

func TestRepository_Close(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}

	repo := &repository{Con: mock}

	mock.ExpectClose()
	if err = repo.Close(ctx); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
