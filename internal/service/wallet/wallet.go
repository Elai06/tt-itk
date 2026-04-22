package wallet

import (
	"context"
	"errors"
	"itk-wallet/internal/dto"
	"itk-wallet/internal/model"
)

//go:generate mockgen -destination=mocks/mock_storage.go -package=mocks itk/internal/service/wallet Storage
type Storage interface {
	Insert(ctx context.Context, wallet model.Wallet) error
	Update(ctx context.Context, wallet model.Wallet) error
	Get(ctx context.Context, uuid int64) (int64, error)
}

type WalletService interface {
	Create(ctx context.Context, received dto.WalletRequest) error
	Get(ctx context.Context, uuid int64) (int64, error)
}

type walletService struct {
	storage Storage
}

func NewWalletService(storage Storage) WalletService {
	return &walletService{storage: storage}
}

func (s *walletService) Create(ctx context.Context, received dto.WalletRequest) error {
	amount, err := s.storage.Get(ctx, received.UUID)

	if err != nil && received.OperationType == model.DEPOSIT {
		return s.storage.Insert(ctx, model.Wallet{
			UUID:    received.UUID,
			Balance: received.Amount,
		})
	}

	switch received.OperationType {
	case model.DEPOSIT:
		return s.storage.Update(ctx, model.Wallet{
			UUID:    received.UUID,
			Balance: amount + received.Amount,
		})
	case model.WITHDRAW:
		if amount >= received.Amount {
			return s.storage.Update(ctx, model.Wallet{
				UUID:    received.UUID,
				Balance: amount - received.Amount,
			})
		}
		return errors.New("not enough wallet")
	}

	return nil
}

func (s *walletService) Get(ctx context.Context, uuid int64) (int64, error) {
	return s.storage.Get(ctx, uuid)
}
