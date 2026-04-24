package server

import (
	"itk-wallet/internal/auth"
	"itk-wallet/internal/handler"
	"itk-wallet/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Run(h handler.Handlers, j auth.Jwt, port string) error {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/user/login", h.Auth.Login)
		v1.POST("/user/register", h.Auth.Register)
		v1.GET("/exchange/rates", h.Exchanger.GetCurrencyRates)
		v1.POST("/exchange", h.Exchanger.ExchangeCurrency)

		protected := v1.Group("")
		protected.Use(middleware.AuthRequired(j))
		{
			protected.POST("/wallet", h.Wallet.Create)
			protected.GET("/wallets/:walletId", h.Wallet.Get)

		}
	}

	return r.Run(":" + port)
}
