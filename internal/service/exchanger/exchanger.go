package exchanger

import "context"

type Storage interface {
	GetCurrencies(ctx context.Context) (map[string]float32, error)
	GetExchangeFromCurrency(ctx context.Context, from, to string) (float32, error)
}

type ExchangerService interface {
	GetCurrencyRates(ctx context.Context) (map[string]float32, error)
	ExchangeCurrency(ctx context.Context, from, to string, amount float32) (float32, error)
}

type exchanger struct {
	storage Storage
}

func NewExchangerService(storage Storage) ExchangerService {
	return &exchanger{
		storage: storage,
	}
}

func (e *exchanger) ExchangeCurrency(ctx context.Context, from, to string, amount float32) (float32, error) {
	rate, err := e.storage.GetExchangeFromCurrency(ctx, from, to)
	if err != nil {
		return 0, err
	}

	return amount / rate, nil
}

func (e *exchanger) GetCurrencyRates(ctx context.Context) (map[string]float32, error) {
	currencies, err := e.storage.GetCurrencies(ctx)
	if err != nil {
		return nil, err
	}

	return currencies, nil
}
