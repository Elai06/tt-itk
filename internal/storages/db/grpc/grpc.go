package grpc

import (
	"context"
	"fmt"
	"itk-wallet/generated/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GprcClient interface {
	Close() error
	GetCurrencyRates(ctx context.Context) (map[string]float32, error)
	GetExchangeForCurrency(ctx context.Context, from, to string) (float32, error)
}

type client struct {
	conn   *grpc.ClientConn
	client exchangerpb.ExchangeServiceClient
}

func NewClient(addr string) (GprcClient, error) {
	fmt.Println(addr)
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &client{
		conn:   conn,
		client: exchangerpb.NewExchangeServiceClient(conn),
	}, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) GetCurrencyRates(ctx context.Context) (map[string]float32, error) {
	resp, err := c.client.GetExchangeRates(ctx, &exchangerpb.Empty{})
	if err != nil {
		return nil, err
	}

	return resp.Rates, nil
}

func (c *client) GetExchangeForCurrency(ctx context.Context, from, to string) (float32, error) {
	resp, err := c.client.GetExchangeRateForCurrency(ctx, &exchangerpb.CurrencyRequest{
		FromCurrency: from,
		ToCurrency:   to,
	})

	if err != nil {
		return 0, err
	}

	return resp.Rate, nil
}
