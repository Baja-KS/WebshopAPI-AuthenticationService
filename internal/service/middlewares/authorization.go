package middlewares

import (
	"context"
	"errors"
	"github.com/Baja-KS/WebshopAPI-AuthenticationService/internal/database"
	"github.com/Baja-KS/WebshopAPI-AuthenticationService/internal/service"
	"gorm.io/gorm"
)

type AuthorizationMiddleware struct {
	DB   *gorm.DB
	Next service.Service
}

func AuthorizeAdmin(user database.UserOut) bool {
	return user.IsAdmin
}

func (a *AuthorizationMiddleware) Login(ctx context.Context, username string, password string) (database.UserOut, string, error) {
	return a.Next.Login(ctx, username, password)
}

func (a *AuthorizationMiddleware) Register(ctx context.Context, user database.UserIn) (string, error) {
	if !database.AuthorizeUser(ctx, a.DB, AuthorizeAdmin) {
		return "Unauthorized", errors.New("unauthorized")
	}
	return a.Next.Register(ctx, user)
}

func (a *AuthorizationMiddleware) GetAll(ctx context.Context) ([]database.UserOut, error) {
	if !database.AuthorizeUser(ctx, a.DB, AuthorizeAdmin) {
		return []database.UserOut{}, errors.New("unauthorized")
	}
	return a.Next.GetAll(ctx)
}

func (a *AuthorizationMiddleware) AuthUser(ctx context.Context) (database.UserOut, error) {
	return a.Next.AuthUser(ctx)
}
