package handler

import (
	"itk-wallet/internal/dto"
	walletService "itk-wallet/internal/service/wallet"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WalletHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
}

type wallet struct {
	service walletService.WalletService
}

func NewWalletHandler(service walletService.WalletService) WalletHandler {
	return &wallet{service: service}
}

func (w *wallet) Create(c *gin.Context) {
	var req dto.WalletRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	err = w.service.Create(c.Request.Context(), req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
	})
}

func (w *wallet) Get(c *gin.Context) {
	uuid, err := strconv.ParseInt(c.Param("walletId"), 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	amount, err := w.service.Get(c.Request.Context(), uuid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, dto.WalletResponse{
		UUID:   uuid,
		Amount: amount,
	})
}
