package handler

import (
	"itk-wallet/internal/dto"
	"itk-wallet/internal/service/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	Logout(c *gin.Context)
}

type auth struct {
	service user.UserService
}

func NewAuthHandler(service user.UserService) AuthHandler {
	return &auth{
		service: service,
	}
}

func (a *auth) Login(c *gin.Context) {
	loginData := dto.LoginRequest{}
	err := c.ShouldBindQuery(&loginData)

	token, err := a.service.Login(c, loginData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (a *auth) Register(c *gin.Context) {
	registerData := dto.RegisterRequest{}
	err := c.ShouldBindQuery(&registerData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	err = a.service.Register(c, registerData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (a *auth) Logout(c *gin.Context) {

}
