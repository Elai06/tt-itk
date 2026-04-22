package wallet

import (
	"context"
	"errors"
	"itk-wallet/internal/model"
	"regexp"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v5"
)

func newWalletRepoMock(t *testing.T) (*Wallet, pgxmock.PgxConnIface) {
	t.Helper()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	return &Wallet{db: mock}, mock
}

func TestWalletRepo_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		wallet    model.Wallet
		mockSetup func(mock pgxmock.PgxConnIface)
		wantErr   bool
	}{
		{
			name:   "success",
			wallet: model.Wallet{UUID: 1, Balance: 100},
			mockSetup: func(mock pgxmock.PgxConnIface) {
				query := regexp.QuoteMeta(`INSERT INTO wallets (balance, uuid)
				VALUES ($1, $2)`)
				mock.ExpectExec(query).
					WithArgs(int64(100), int64(1)).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
		},
		{
			name:   "db error",
			wallet: model.Wallet{UUID: 2, Balance: 200},
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := newWalletRepoMock(t)
			tt.mockSetup(mock)

			err := repo.Insert(ctx, tt.wallet)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

func TestWalletRepo_Update(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	selectQuery := regexp.QuoteMeta(`SELECT balance FROM wallets WHERE uuid = $1 FOR UPDATE`)
	updateQuery := regexp.QuoteMeta(`
					UPDATE  wallets
					SET balance = $1
					WHERE uuid = $2
		`)

	tests := []struct {
		name      string
		wallet    model.Wallet
		mockSetup func(mock pgxmock.PgxConnIface)
		wantErr   bool
	}{
		{
			name:   "success",
			wallet: model.Wallet{UUID: 1, Balance: 500},
			mockSetup: func(mock pgxmock.PgxConnIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(selectQuery).
					WithArgs(int64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"balance"}).AddRow(int64(400)))
				mock.ExpectExec(updateQuery).
					WithArgs(int64(500), int64(1)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit()
			},
		},
		{
			name:   "db error on exec",
			wallet: model.Wallet{UUID: 2, Balance: 700},
			mockSetup: func(mock pgxmock.PgxConnIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(selectQuery).
					WithArgs(int64(2)).
					WillReturnRows(pgxmock.NewRows([]string{"balance"}).AddRow(int64(600)))
				mock.ExpectExec(updateQuery).
					WithArgs(int64(700), int64(2)).
					WillReturnError(errors.New("update failed"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := newWalletRepoMock(t)
			tt.mockSetup(mock)

			err := repo.Update(ctx, tt.wallet)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

func TestWalletRepo_Get(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	query := regexp.QuoteMeta(`SELECT balance
			FROM wallets
			WHERE uuid = $1`)

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
				mock.ExpectQuery(query).
					WithArgs(int64(1)).
					WillReturnRows(pgxmock.NewRows([]string{"balance"}).AddRow(int64(1000)))
			},
			wantAmount: 1000,
		},
		{
			name: "not found",
			uuid: 2,
			mockSetup: func(mock pgxmock.PgxConnIface) {
				mock.ExpectQuery(query).
					WithArgs(int64(2)).
					WillReturnError(pgx.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "scan error",
			uuid: 3,
			mockSetup: func(mock pgxmock.PgxConnIface) {
				mock.ExpectQuery(query).
					WithArgs(int64(3)).
					WillReturnRows(pgxmock.NewRows([]string{"balance"}).AddRow("bad_value"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, mock := newWalletRepoMock(t)
			tt.mockSetup(mock)

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
