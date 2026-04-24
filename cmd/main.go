package main

import (
	"context"
	"flag"
	"itk-wallet/internal/auth"
	"itk-wallet/internal/config"
	"itk-wallet/internal/handler"
	"itk-wallet/internal/server"
	"itk-wallet/internal/service/exchanger"
	user2 "itk-wallet/internal/service/user"
	"itk-wallet/internal/service/wallet"
	"itk-wallet/internal/storages/db/grpc"
	exchanger2 "itk-wallet/internal/storages/db/grpc/exchanger"
	"itk-wallet/internal/storages/db/postgres"
	"itk-wallet/internal/storages/db/postgres/user"
	walletRepo "itk-wallet/internal/storages/db/postgres/wallet"
	"log"
)

func main() {
	configPath := flag.String("c", "config.env", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	ctx := context.Background()

	client, err := postgres.New(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	defer func() {
		if err = client.Close(ctx); err != nil {
			log.Printf("Error closing db client: %v", err)
		}
	}()

	walletStorage := walletRepo.NewWallet(client.Conn)
	walletSvc := wallet.NewWalletService(walletStorage)

	userStorage := user.NewUser(client.Conn)
	jwt := auth.NewJwt([]byte(cfg.JWTSecret), cfg.JWTExpiration)
	userSvc := user2.NewUserService(userStorage, jwt)

	grpcClient, err := grpc.NewClient(cfg.GrpcAddr)
	if err != nil {
		log.Fatalf("Error connecting to grpc: %v", err)
	}

	exStorage := exchanger2.NewExchanger(grpcClient)
	exSvc := exchanger.NewExchangerService(exStorage)

	walletHandler := handler.NewWalletHandler(walletSvc)
	authHandler := handler.NewAuthHandler(userSvc)
	exHandler := handler.NewExchangerHandler(exSvc)

	handlers := handler.NewHandlers(&walletHandler, &authHandler, &exHandler)

	if err = server.Run(*handlers, jwt, cfg.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
