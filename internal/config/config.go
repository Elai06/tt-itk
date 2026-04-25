package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl         string
	MigrationsDir string
	Port          string
	GrpcAddr      string
	JWTSecret     string
	JWTExpiration time.Duration
}

type KafkaConfig struct {
	Broker  []string
	Topic   string
	GroupId string
}

func Load(path string) (*Config, error) {
	if err := godotenv.Overload(path); err != nil {
		return nil, fmt.Errorf("loading config file %s: %w", path, err)
	}

	dbUrl, err := requiredEnv("DB_URI")
	if err != nil {
		return nil, err
	}

	migrationsDir, err := requiredEnv("MIGRATIONS_DIR")
	if err != nil {
		return nil, err
	}

	port, err := requiredEnv("PORT")
	if err != nil {
		return nil, err
	}

	jwtSecret, err := requiredEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}

	jwtExpiration, err := requiredEnv("JWT_EXPIRED")
	if err != nil {
		return nil, err
	}

	expiration, err := strconv.Atoi(jwtExpiration)
	if err != nil {
		return nil, err
	}

	grpcAddr, err := requiredEnv("GRPC_ADDR")
	if err != nil {
		return nil, err
	}

	return &Config{
		DBUrl:         dbUrl,
		MigrationsDir: migrationsDir,
		Port:          port,
		JWTSecret:     jwtSecret,
		JWTExpiration: time.Duration(expiration) * time.Second,
		GrpcAddr:      grpcAddr,
	}, nil
}

func LoadKafkaConfig(path string) (*KafkaConfig, error) {
	if err := godotenv.Overload(path); err != nil {
		return nil, fmt.Errorf("loading config file %s: %w", path, err)
	}

	broker, err := requiredEnv("KAFKA_BROKER")
	if err != nil {
		return nil, err
	}
	brokers := make([]string, 0)
	brokers = append(brokers, broker)

	topic, err := requiredEnv("KAFKA_TOPIC")
	if err != nil {
		return nil, err
	}

	groupId, err := requiredEnv("KAFKA_GROUP_ID")
	if err != nil {
		return nil, err
	}

	return &KafkaConfig{
		Broker:  brokers,
		Topic:   topic,
		GroupId: groupId,
	}, nil
}

func requiredEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return "", fmt.Errorf("%s is not set", key)
	}
	return value, nil
}
