package main

import (
	"context"
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
	err := godotenv.Load("config.env")

	if err != nil {
		log.Printf("Error loading .env file")
		return
	}

	dbUrl := os.Getenv("DB_URI")
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

	migratorDir := os.Getenv("MIGRATIONS_DIR")
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*rep.GetConConfig()), migratorDir)

	err = migratorRunner.Up()
	if err != nil {
		log.Printf("Error running migrations: %v", err)
		return
	}

	s := service.NewWalletService(rep)
	h := handler.NewWalletHandler(s)
	err = h.RegisterRoutes(os.Getenv("PORT"))
	if err != nil {
		log.Printf("Error registering routes: %v", err)
		return
	}
}
