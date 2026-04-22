package wallet

import (
	"context"
	"errors"
	"itk-wallet/internal/dto"
	"itk-wallet/internal/model"
	"itk-wallet/internal/service/wallet/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestWalletServiceCreate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		req       dto.WalletRequest
		mockSetup func(ctx context.Context, repo *mocks.MockStorage)
		wantErr   string
	}{
		{
			name: "create new wallet on deposit when wallet does not exist",
			req: dto.WalletRequest{
				UUID:          1,
				OperationType: model.DEPOSIT,
				Amount:        100,
			},
			mockSetup: func(ctx context.Context, repo *mocks.MockStorage) {
				gomock.InOrder(
					repo.EXPECT().
						Get(ctx, int64(1)).
						Return(int64(0), errors.New("not found")),
					repo.EXPECT().
						Create(ctx, model.Wallet{UUID: 1, Balance: 100}).
						Return(nil),
				)
			},
		},
		{
			name: "increase balance on deposit",
			req: dto.WalletRequest{
				UUID:          2,
				OperationType: model.DEPOSIT,
				Amount:        50,
			},
			mockSetup: func(ctx context.Context, repo *mocks.MockStorage) {
				gomock.InOrder(
					repo.EXPECT().
						Get(ctx, int64(2)).
						Return(int64(150), nil),
					repo.EXPECT().
						Update(ctx, model.Wallet{UUID: 2, Balance: 200}).
						Return(nil),
				)
			},
		},
		{
			name: "decrease balance on withdraw",
			req: dto.WalletRequest{
				UUID:          3,
				OperationType: model.WITHDRAW,
				Amount:        30,
			},
			mockSetup: func(ctx context.Context, repo *mocks.MockStorage) {
				gomock.InOrder(
					repo.EXPECT().
						Get(ctx, int64(3)).
						Return(int64(100), nil),
					repo.EXPECT().
						Update(ctx, model.Wallet{UUID: 3, Balance: 70}).
						Return(nil),
				)
			},
		},
		{
			name: "return error when balance is insufficient",
			req: dto.WalletRequest{
				UUID:          4,
				OperationType: model.WITHDRAW,
				Amount:        80,
			},
			mockSetup: func(ctx context.Context, repo *mocks.MockStorage) {
				repo.EXPECT().
					Get(ctx, int64(4)).
					Return(int64(50), nil)
			},
			wantErr: "not enough wallet",
		},
		{
			name: "return create error",
			req: dto.WalletRequest{
				UUID:          5,
				OperationType: model.DEPOSIT,
				Amount:        100,
			},
			mockSetup: func(ctx context.Context, repo *mocks.MockStorage) {
				createErr := errors.New("create failed")
				gomock.InOrder(
					repo.EXPECT().
						Get(ctx, int64(5)).
						Return(int64(0), errors.New("not found")),
					repo.EXPECT().
						Create(ctx, model.Wallet{UUID: 5, Balance: 100}).
						Return(createErr),
				)
			},
			wantErr: "create failed",
		},
		{
			name: "return update error on deposit",
			req: dto.WalletRequest{
				UUID:          6,
				OperationType: model.DEPOSIT,
				Amount:        20,
			},
			mockSetup: func(ctx context.Context, repo *mocks.MockStorage) {
				updateErr := errors.New("update failed")
				gomock.InOrder(
					repo.EXPECT().
						Get(ctx, int64(6)).
						Return(int64(10), nil),
					repo.EXPECT().
						Update(ctx, model.Wallet{UUID: 6, Balance: 30}).
						Return(updateErr),
				)
			},
			wantErr: "update failed",
		},
		{
			name: "return update error on withdraw",
			req: dto.WalletRequest{
				UUID:          7,
				OperationType: model.WITHDRAW,
				Amount:        10,
			},
			mockSetup: func(ctx context.Context, repo *mocks.MockStorage) {
				updateErr := errors.New("update failed")
				gomock.InOrder(
					repo.EXPECT().
						Get(ctx, int64(7)).
						Return(int64(15), nil),
					repo.EXPECT().
						Update(ctx, model.Wallet{UUID: 7, Balance: 5}).
						Return(updateErr),
				)
			},
			wantErr: "update failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mocks.NewMockStorage(ctrl)
			svc := NewWalletService(repo)
			ctx := context.Background()

			tt.mockSetup(ctx, repo)

			err := svc.Create(ctx, tt.req)
			if tt.wantErr == "" && err != nil {
				t.Fatalf("Insert() unexpected error = %v", err)
			}
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("Insert() error = nil, want %q", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("Insert() error = %q, want %q", err.Error(), tt.wantErr)
				}
			}
		})
	}
}

func TestWalletServiceGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		uuid       int64
		mockSetup  func(ctx context.Context, repo *mocks.MockStorage)
		wantAmount int64
		wantErr    string
	}{
		{
			name: "success",
			uuid: 42,
			mockSetup: func(ctx context.Context, repo *mocks.MockStorage) {
				repo.EXPECT().
					Get(ctx, int64(42)).
					Return(int64(900), nil)
			},
			wantAmount: 900,
		},
		{
			name: "storage error",
			uuid: 77,
			mockSetup: func(ctx context.Context, repo *mocks.MockStorage) {
				repo.EXPECT().
					Get(ctx, int64(77)).
					Return(int64(0), errors.New("not found"))
			},
			wantErr: "not found",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mocks.NewMockStorage(ctrl)
			svc := NewWalletService(repo)
			ctx := context.Background()

			tt.mockSetup(ctx, repo)

			got, err := svc.Get(ctx, tt.uuid)
			if tt.wantErr == "" && err != nil {
				t.Fatalf("Get() unexpected error = %v", err)
			}
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("Get() error = nil, want %q", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("Get() error = %q, want %q", err.Error(), tt.wantErr)
				}
			}
			if got != tt.wantAmount {
				t.Fatalf("Get() = %d, want %d", got, tt.wantAmount)
			}
		})
	}
}
