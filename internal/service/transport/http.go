package transport

import (
	"context"
	"errors"
	"github.com/Baja-KS/Webshop-API/AuthenticationService/internal/service/endpoints"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"
)

func GetAuthToken(r *http.Request) (string,error) {
	authHeader:=r.Header["Authorization"]
	if len(authHeader)==0 {
		return "", errors.New("no auth header")
	}
	authHeaderParts:=strings.Split(authHeader[0]," ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "",errors.New("invalid auth header")
	}
	return authHeaderParts[1],nil
}

func AuthExtractor(ctx context.Context, r *http.Request) context.Context {
	token,err:=GetAuthToken(r)
	if err != nil {
		return context.WithValue(ctx,"auth","")
	}
	return context.WithValue(ctx,"auth",token)
}

func NewHTTPHandler(ep endpoints.EndpointSet) http.Handler {
	router:=mux.NewRouter()

	loginHandler:=httptransport.NewServer(ep.LoginEndpoint,endpoints.DecodeLoginRequest,endpoints.EncodeResponse,httptransport.ServerBefore(AuthExtractor))
	registerHandler:=httptransport.NewServer(ep.RegisterEndpoint,endpoints.DecodeRegisterRequest,endpoints.EncodeResponse,httptransport.ServerBefore(AuthExtractor))
	getAllHandler:=httptransport.NewServer(ep.GetAllEndpoint,endpoints.DecodeGetAllRequest,endpoints.EncodeResponse,httptransport.ServerBefore(AuthExtractor))
	authUserHandler:=httptransport.NewServer(ep.AuthUserEndpoint,endpoints.DecodeAuthUserRequest,endpoints.EncodeResponse,httptransport.ServerBefore(AuthExtractor))
	router.Handle("/Login",loginHandler).Methods(http.MethodPost)
	router.Handle("/Register",registerHandler).Methods(http.MethodPost)
	router.Handle("/GetAll",getAllHandler).Methods(http.MethodGet)
	router.Handle("/User",authUserHandler).Methods(http.MethodGet)
	router.Handle("/metrics",promhttp.Handler()).Methods(http.MethodGet)

	return router
}