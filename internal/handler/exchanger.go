package handler

import (
	"itk-wallet/internal/service/exchanger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ExchangerHandler interface {
	GetCurrencyRates(c *gin.Context)
	ExchangeCurrency(c *gin.Context)
}

type exchangerHandler struct {
	service exchanger.ExchangerService
}

func NewExchangerHandler(service exchanger.ExchangerService) ExchangerHandler {
	return &exchangerHandler{
		service: service,
	}
}

func (e *exchangerHandler) GetCurrencyRates(c *gin.Context) {
	rates, err := e.service.GetCurrencyRates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rates": rates,
	})
}

func (e *exchangerHandler) ExchangeCurrency(c *gin.Context) {
	from := c.Query("from_currency")
	to := c.Query("to_currency")
	amount := c.Query("amount")

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	rate, err := e.service.ExchangeCurrency(c.Request.Context(), from, to, float32(amountFloat))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exchange successful",
		from:      amount,
		to:        rate,
	})
}
