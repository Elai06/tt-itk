package wallet

import (
	"context"
	"errors"
	"fmt"
	"itk-wallet/internal/dto"
	"itk-wallet/internal/model"
	"itk-wallet/internal/storages/kafka"
	"log"
)

const TRIGGER_EVENT int64 = 30000

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
	storage   Storage
	kProducer kafka.Producer
}

func NewWalletService(storage Storage, kProducer kafka.Producer) WalletService {
	return &walletService{storage: storage, kProducer: kProducer}
}

func (s *walletService) Create(ctx context.Context, received dto.WalletRequest) error {
	amount, err := s.storage.Get(ctx, received.UUID)

	if err != nil && received.OperationType == model.DEPOSIT {
		err = s.storage.Insert(ctx, model.Wallet{
			UUID:         received.UUID,
			CurrencyCode: received.CurrencyCode,
			Balance:      received.Amount,
		})

		if err != nil {
			return err
		}

		// похорошему вынести в конфиг
		if received.Amount > TRIGGER_EVENT {
			s.sendEvent(ctx, received.CurrencyCode, received.UUID, received.Amount+amount)
		}

		return nil
	}

	switch received.OperationType {
	case model.DEPOSIT:
		err = s.storage.Update(ctx, model.Wallet{
			UUID:         received.UUID,
			CurrencyCode: received.CurrencyCode,
			Balance:      amount + received.Amount,
		})

		if err != nil {
			return err
		}

		// похорошему вынести в конфиг
		if received.Amount > TRIGGER_EVENT {
			s.sendEvent(ctx, received.CurrencyCode, received.UUID, received.Amount+amount)
		}

		return nil
	case model.WITHDRAW:
		if amount >= received.Amount {
			return s.storage.Update(ctx, model.Wallet{
				UUID:         received.UUID,
				CurrencyCode: received.CurrencyCode,
				Balance:      amount - received.Amount,
			})
		}
		return errors.New("not enough wallet")
	}

	return nil
}

func (s *walletService) Get(ctx context.Context, uuid int64) (int64, error) {
	return s.storage.Get(ctx, uuid)
}

func (s *walletService) sendEvent(ctx context.Context, code string, id, amount int64) {
	value := kafka.WalletEventProducer{
		UUID:         id,
		CurrencyCode: code,
		Amount:       amount,
	}
	key := fmt.Sprint("walletID-", id)
	err := s.kProducer.SendRemittanceWallet(ctx, key, value)
	if err != nil {
		log.Printf("failed to send remittance wallet event to kafka: %v", err)
	}
}
