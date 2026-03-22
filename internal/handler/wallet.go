package handler

import (
	"itk/internal/dto"
	"itk/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WalletHandler interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	RegisterRoutes(port string) error
}

type wallet struct {
	service service.WalletService
}

func NewWalletHandler(service service.WalletService) WalletHandler {
	return &wallet{service: service}
}

func (w *wallet) RegisterRoutes(port string) error {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/wallet", w.Create)
		v1.GET("/wallets/:walletId", w.Get)
	}

	err := r.Run(":" + port)
	if err != nil {
		return err
	}

	return nil
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

	return
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
