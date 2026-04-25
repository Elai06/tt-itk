package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"itk-wallet/internal/config"
	"time"

	"github.com/segmentio/kafka-go"
)

type WalletEventProducer struct {
	UUID         int64
	CurrencyCode string
	Amount       int64
}

type Producer interface {
	SendRemittanceWallet(ctx context.Context, key string, value WalletEventProducer) error
	Close() error
}

type producer struct {
	writer *kafka.Writer
}

func NewProduce(cfg config.KafkaConfig) Producer {
	return &producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Broker...),
			Topic:        cfg.Topic,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireAll,
			BatchTimeout: 10 * time.Millisecond,
		},
	}
}

func (p *producer) SendRemittanceWallet(ctx context.Context, key string, value WalletEventProducer) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: payload,
		Time:  time.Now(),
	})

	if err != nil {
		return fmt.Errorf("kafka: write messages: %w", err)
	}

	return nil
}

func (p *producer) Close() error {
	return p.writer.Close()
}
