package auth

import (
	"fmt"
	"sync"
	"time"

	jwt2 "github.com/golang-jwt/jwt/v5"
)

type Jwt interface {
	GenerateToken(userId int64, username, email string) (string, error)
	CheckToken(tokenString string) error
}

type jwt struct {
	secret []byte
	t      time.Duration
	tokens map[int64]string
	mu     sync.RWMutex
}

type Claims struct {
	Username string
	Email    string

	jwt2.RegisteredClaims
}

func NewJwt(secret []byte, time time.Duration) Jwt {
	return &jwt{
		secret: secret,
		t:      time,
		tokens: make(map[int64]string),
		mu:     sync.RWMutex{},
	}
}

func (j *jwt) GenerateToken(userId int64, username, email string) (string, error) {
	claims := Claims{
		Username: username,
		Email:    email,

		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(time.Now().Add(j.t)),
			IssuedAt:  jwt2.NewNumericDate(time.Now()),
		},
	}

	token := jwt2.NewWithClaims(jwt2.SigningMethodHS256, claims)
	strToken, err := token.SignedString(j.secret)

	if err != nil {
		return "", err
	}

	j.mu.RLock()
	defer j.mu.RUnlock()
	j.tokens[userId] = strToken

	return strToken, nil
}

func (j *jwt) CheckToken(tokenString string) error {
	claims := Claims{}

	keyfunc := func(token *jwt2.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt2.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return j.secret, nil
	}

	token, err := jwt2.ParseWithClaims(tokenString, &claims, keyfunc)
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
