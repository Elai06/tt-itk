package kafka

import (
	"itk-wallet/internal/config"
)

type kafkaClient struct {
	producer Producer
}

type Kafka interface {
	Producer() Producer
	Close() error
}

func NewKafka(cfg config.KafkaConfig) (Kafka, error) {
	prod := NewProduce(cfg)
	return &kafkaClient{producer: prod}, nil
}

func (k *kafkaClient) Producer() Producer {
	return k.producer
}
func (k *kafkaClient) Close() error {
	err := k.producer.Close()
	return err
}
