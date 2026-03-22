package main

import (
	"context"
	"fmt"
	"itk/internal/handler"
	"itk/internal/migrator"
	"itk/internal/repository"
	"itk/internal/service"
	"log"
	"os"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Overload("config.env")

	if err != nil {
		log.Printf("Error loading config.env: %v", err)
		return
	}

	dbUrl, err := requiredEnv("DB_URI")
	if err != nil {
		log.Printf("Configuration error: %v", err)
		return
	}

	ctx := context.Background()
	rep, err := repository.NewRepository(dbUrl, ctx)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return
	}

	defer func(rep repository.Repository) {
		err = rep.Close(ctx)
		if err != nil {
			log.Printf("Error closing repository: %v", err)
			return
		}
	}(rep)

	migratorDir, err := requiredEnv("MIGRATIONS_DIR")
	if err != nil {
		log.Printf("Configuration error: %v", err)
		return
	}

	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*rep.GetConConfig()), migratorDir)

	err = migratorRunner.Up()
	if err != nil {
		log.Printf("Error running migrations: %v", err)
		return
	}

	s := service.NewWalletService(rep)
	h := handler.NewWalletHandler(s)
	port, err := requiredEnv("PORT")
	if err != nil {
		log.Printf("Configuration error: %v", err)
		return
	}

	err = h.RegisterRoutes(port)
	if err != nil {
		log.Printf("Error registering routes: %v", err)
		return
	}
}

func requiredEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return "", fmt.Errorf("%s is not set in config.env", key)
	}

	return value, nil
}
