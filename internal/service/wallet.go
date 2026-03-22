package service

import (
	"context"
	"errors"
	"itk/internal/dto"
	"itk/internal/model"
	"itk/internal/repository"
)

//go:generate mockgen -destination=mocks/mock_wallet_service.go -package=mocks itk/internal/service WalletService

type WalletService interface {
	Create(ctx context.Context, received dto.WalletRequest) error
	Get(ctx context.Context, uuid int64) (int64, error)
}

type walletService struct {
	repo repository.Repository
}

func NewWalletService(repo repository.Repository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) Create(ctx context.Context, received dto.WalletRequest) error {
	amount, err := s.repo.Get(ctx, received.UUID)
	if err != nil && received.OperationType == model.DEPOSIT {
		return s.repo.Create(ctx, received)
	} else {
		switch received.OperationType {
		case model.DEPOSIT:
			received.Amount = amount + received.Amount
			return s.repo.Update(ctx, received)
		case model.WITHDRAW:
			isAvailable := amount >= received.Amount
			if isAvailable {
				received.Amount = amount - received.Amount
				return s.repo.Update(ctx, received)
			} else {
				return errors.New("not enough wallet")
			}
		}

		return nil
	}
}

func (s *walletService) Get(ctx context.Context, uuid int64) (int64, error) {
	return s.repo.Get(ctx, uuid)
}
