package database

import (
	"context"
	"errors"
	stdjwt "github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

type User struct {
	gorm.Model
	ID uint `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Username string `gorm:"not null;unique" json:"username"`
	Fullname string `gorm:"not null" json:"fullname"`
	Email string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null"  json:"password"`
	IsAdmin bool `gorm:"default:false" json:"isAdmin"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt,omitempty"`
}

func (u *User) Out() UserOut {
	if u.Username=="" {
		var emptyUser UserOut
		return emptyUser
	}
	userOut:=UserOut{
		ID:        u.ID,
		Username:  u.Username,
		Fullname:  u.Fullname,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	return userOut
}

type UserIn struct {
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Email string `json:"email"`
	Password string `json:"password"`
	IsAdmin bool `json:"isAdmin"`
}

type UserOut struct {
	ID uint `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Email string `json:"email"`
	IsAdmin bool `json:"isAdmin"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

func UserArrayOut(models []User) []UserOut {
	outArr:=make([]UserOut,len(models))
	for i,item := range models {
		outArr[i]=item.Out()
	}
	return outArr
}

func GetAuthUser(ctx context.Context, db *gorm.DB) (UserOut, error) {
	tokenString:=ctx.Value("auth").(string)
	parsed,err:=stdjwt.ParseWithClaims(tokenString,&stdjwt.StandardClaims{}, func(token *stdjwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")),nil
	})
	if err != nil {
		return UserOut{}, err
	}
	if parsed.Valid {
		castedClaims:=parsed.Claims.(*stdjwt.StandardClaims)
		userId:=castedClaims.Id
		id,err:=strconv.ParseUint(userId,10,64)
		if err != nil {
			return UserOut{}, err
		}
		var userDB User
		notFound:=db.Where("id = ?",id).First(&userDB).Error
		if notFound != nil {
			return UserOut{}, notFound
		}
		return userDB.Out(), nil
	}
	return UserOut{}, errors.New("Not found")
}

type AuthorizationFunction func(UserOut) bool

func AuthenticateUser(ctx context.Context, db *gorm.DB) bool {
	return AuthorizeUser(ctx,db, func(user UserOut) bool {
		return user.ID!=0
	})
}

func AuthorizeUser(ctx context.Context, db *gorm.DB, fn AuthorizationFunction) bool {
	user,err:=GetAuthUser(ctx,db)
	if err != nil {
		return false
	}
	return fn(user)
}