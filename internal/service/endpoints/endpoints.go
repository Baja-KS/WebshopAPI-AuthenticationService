package endpoints

import (
	"context"
	"github.com/Baja-KS/Webshop-API/AuthenticationService/internal/database"
	"github.com/Baja-KS/Webshop-API/AuthenticationService/internal/service"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/crypto/bcrypt"
)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type EndpointSet struct {
	LoginEndpoint endpoint.Endpoint
	RegisterEndpoint endpoint.Endpoint
	GetAllEndpoint endpoint.Endpoint
	AuthUserEndpoint endpoint.Endpoint
}

func NewEndpointSet(svc service.Service) EndpointSet {
	//var authUserEndpoint endpoint.Endpoint
	//authUserEndpoint=MakeAuthUserEndpoint(svc)
	//authUserEndpoint=middlewares.Authenticated()(authUserEndpoint)
	//var authUserEndpoint endpoint.Endpoint
	//	kf:=func(token *stdjwt.Token) (interface{}, error) { return []byte(os.Getenv("JWT_KEY")), nil }
	//	authUserEndpoint=MakeAuthUserEndpoint(svc)
	//	authUserEndpoint=jwt.NewParser(kf,stdjwt.SigningMethodHS256,jwt.StandardClaimsFactory)(authUserEndpoint)
	return EndpointSet{
		LoginEndpoint:    MakeLoginEndpoint(svc),
		RegisterEndpoint: MakeRegisterEndpoint(svc),
		GetAllEndpoint:   MakeGetAllEndpoint(svc),
		AuthUserEndpoint: MakeAuthUserEndpoint(svc),
	}
}

func MakeLoginEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req:=request.(LoginRequest)
		user,token,e:=svc.Login(ctx,req.Username,req.Password)
		if e != nil {
			return LoginResponse{Message: "Wrong username or password"},nil
		}
		userOut:=database.UserOut{
			ID:        user.ID,
			Username:  user.Username,
			Fullname:  user.Fullname,
			Email:     user.Email,
			IsAdmin:   user.IsAdmin,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		return LoginResponse{User: userOut,Token: token,Message: "Login successful"},nil
	}
}

func MakeRegisterEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req:=request.(RegisterRequest)

		user:=database.UserIn{
			Username:  req.Username,
			Fullname:  req.Fullname,
			Email:    req.Email ,
			Password:  req.Password,
			IsAdmin:   req.IsAdmin,
		}
		msg,err:=svc.Register(ctx,user)
		if err != nil {
			return RegisterResponse{Message: err.Error()},err
		}
		return RegisterResponse{Username: user.Username,Message: msg},nil
	}
}

func MakeGetAllEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		users,err:=svc.GetAll(ctx)
		if err != nil {
			var empty []database.UserOut=nil
			return GetAllResponse{Users: empty},err
		}
		var usersOut []database.UserOut
		for _, user := range users {
			usersOut=append(usersOut,database.UserOut{
				ID:       user.ID,
				Username: user.Username,
				Fullname: user.Fullname,
				Email:    user.Email,
				IsAdmin:  user.IsAdmin,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			})
		}
		return GetAllResponse{Users: usersOut},nil
	}
}

func MakeAuthUserEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user,err:=svc.AuthUser(ctx)
		if err != nil {
			return AuthUserResponse{User: user},err
		}
		return AuthUserResponse{User: user},nil
	}
}
