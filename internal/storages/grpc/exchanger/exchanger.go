package exchanger

import (
	"context"
	exchanger2 "itk-wallet/internal/service/exchanger"
	"itk-wallet/internal/storages/grpc"
)

type exchanger struct {
	client grpc.GprcClient
}

func NewExchanger(client grpc.GprcClient) exchanger2.Storage {
	return &exchanger{client: client}
}

func (e *exchanger) GetCurrencies(ctx context.Context) (map[string]float32, error) {
	rates, err := e.client.GetCurrencyRates(ctx)
	if err != nil {
		return nil, err
	}

	return rates, nil
}

func (e *exchanger) GetExchangeFromCurrency(ctx context.Context, from, to string) (float32, error) {
	rate, err := e.client.GetExchangeForCurrency(ctx, from, to)
	if err != nil {
		return 0, err
	}

	return rate, nil
}
