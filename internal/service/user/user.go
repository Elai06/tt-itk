package user

import (
	"context"
	"itk-wallet/internal/auth"
	"itk-wallet/internal/dto"
	"itk-wallet/internal/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -destination=mocks/mock_storage.go -package=mocks itk/internal/service/user Storage
type Storage interface {
	Insert(ctx context.Context, user model.User) error
	Get(ctx context.Context, email string) (model.User, error)
}

type UserService interface {
	Login(ctx *gin.Context, l dto.LoginRequest) (string, error)
	Register(c *gin.Context, r dto.RegisterRequest) error
}

type userService struct {
	storage Storage
	jwt     auth.Jwt
}

func NewUserService(storage Storage, jwt auth.Jwt) UserService {
	return &userService{
		storage: storage,
		jwt:     jwt,
	}
}

func (u *userService) Login(ctx *gin.Context, l dto.LoginRequest) (string, error) {
	usrModel, err := u.storage.Get(ctx, l.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usrModel.Password), []byte(l.Password))
	if err != nil {
		return "", err
	}

	token, err := u.jwt.GenerateToken(usrModel.ID, usrModel.Username, usrModel.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *userService) Register(ctx *gin.Context, r dto.RegisterRequest) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	usrModel := model.User{
		Email:    r.Email,
		Password: string(hashPassword),
		Username: r.Username,
	}

	err = u.storage.Insert(ctx, usrModel)
	if err != nil {
		return err
	}

	return nil
}
