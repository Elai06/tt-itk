package kafka

import (
	"itk-wallet/internal/config"
)

type kafkaClient struct {
	wallet   WalletProducer
	analytic AnalyticProducer
}

type Kafka interface {
	WalletProducer() WalletProducer
	AnalyticProducer() AnalyticProducer
	Close() error
}

func NewKafka(cfg config.KafkaConfig) (Kafka, error) {
	w := NewWalletProducer(cfg)
	a := NewAnalyticProducer(cfg)
	return &kafkaClient{wallet: w, analytic: a}, nil
}

func (k *kafkaClient) WalletProducer() WalletProducer {
	return k.wallet
}

func (k *kafkaClient) AnalyticProducer() AnalyticProducer {
	return k.analytic
}

func (k *kafkaClient) Close() error {
	err := k.wallet.Close()
	return err
}
