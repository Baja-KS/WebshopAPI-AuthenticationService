package service

import (
	"context"
	"errors"
	"github.com/Baja-KS/WebshopAPI-AuthenticationService/internal/database"
	stdjwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
	"strconv"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type Service interface {
	Login(ctx context.Context, username string, password string) (database.UserOut, string, error)
	Register(ctx context.Context, user database.UserIn) (string, error)
	GetAll(ctx context.Context) ([]database.UserOut, error)
	AuthUser(ctx context.Context) (database.UserOut, error)
}

func (a *AuthenticationService) Login(ctx context.Context, username string, password string) (database.UserOut, string, error) {
	var userDB database.User
	err := a.DB.Where("username = ?", username).First(&userDB).Error
	if err != nil {
		return database.UserOut{}, "", err
	}
	if !CheckPasswordHash(password, userDB.Password) {
		return database.UserOut{}, "", errors.New("Login failed")
	}
	at := stdjwt.NewWithClaims(stdjwt.SigningMethodHS256, stdjwt.StandardClaims{
		Id: strconv.Itoa(int(userDB.ID)),
	})
	token, err := at.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return database.UserOut{}, "", err
	}
	return userDB.Out(), token, nil
}

func (a *AuthenticationService) Register(ctx context.Context, user database.UserIn) (string, error) {
	hash, e := HashPassword(user.Password)
	if e != nil {
		return "Error", e
	}
	dbUser := database.User{
		Username: user.Username,
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: hash,
		IsAdmin:  false,
	}
	result := a.DB.Create(&dbUser)
	if result.Error != nil {
		return "Error", result.Error
	}
	return "Register successful", nil
}

func (a *AuthenticationService) GetAll(ctx context.Context) ([]database.UserOut, error) {

	var users []database.User
	result := a.DB.Find(&users)

	if result.Error != nil {
		return database.UserArrayOut(users), result.Error
	}
	return database.UserArrayOut(users), nil
}

func (a *AuthenticationService) AuthUser(ctx context.Context) (database.UserOut, error) {
	return database.GetAuthUser(ctx, a.DB)
}

type AuthenticationService struct {
	DB *gorm.DB
}
