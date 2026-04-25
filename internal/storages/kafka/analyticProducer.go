package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"itk-wallet/internal/config"
	"time"

	"github.com/segmentio/kafka-go"
)

type AnalyticEventProducer struct {
	UUID         int64
	CurrencyCode string
	Amount       int64
}

type AnalyticProducer interface {
	SendEvent(ctx context.Context, key string, value AnalyticEventProducer) error
	Close() error
}

type analytic struct {
	writer *kafka.Writer
}

func NewAnalyticProducer(cfg config.KafkaConfig) AnalyticProducer {
	return &analytic{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Broker...),
			Topic:        cfg.AnalyticTopic,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireAll,
			BatchTimeout: 10 * time.Millisecond,
		},
	}
}

func (p *analytic) SendEvent(ctx context.Context, key string, data AnalyticEventProducer) error {
	payload, err := json.Marshal(data)
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

func (p *analytic) Close() error {
	return p.writer.Close()
}
