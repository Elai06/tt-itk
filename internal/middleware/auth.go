package middleware

import (
	"errors"
	"itk-wallet/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwt2 "github.com/golang-jwt/jwt/v5"
)

func AuthRequired(j auth.Jwt) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization format"})
			return
		}

		tokenString := parts[1]
		err := j.CheckToken(tokenString)
		if err != nil {
			status := http.StatusUnauthorized
			msg := "invalid token"
			if errors.Is(err, jwt2.ErrTokenExpired) {
				msg = "token expired"
			}
			c.AbortWithStatusJSON(status, gin.H{"error": msg})
			return
		}

		c.Next()
	}
}
