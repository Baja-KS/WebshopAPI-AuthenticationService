package endpoints

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Baja-KS/WebshopAPI-AuthenticationService/internal/database"
	"net/http"
	"strings"
)

func GetAuthToken(r *http.Request) (string, error) {
	authHeader := r.Header["Authorization"]
	if len(authHeader) == 0 {
		return "", errors.New("no auth header")
	}
	authHeaderParts := strings.Split(authHeader[0], " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("invalid auth header")
	}
	return authHeaderParts[1], nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User    database.UserOut `json:"user,omitempty"`
	Token   string           `json:"token,omitempty"`
	Message string           `json:"message"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"isAdmin"`
}

type RegisterResponse struct {
	//ID uint `json:"id"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

type GetAllRequest struct {
}

type GetAllResponse struct {
	Users []database.UserOut `json:"users"`
}

type AuthUserRequest struct {
}

type AuthUserResponse struct {
	User database.UserOut `json:"user,omitempty"`
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return json.NewEncoder(w).Encode(response)
}

func DecodeLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeRegisterRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeGetAllRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request GetAllRequest
	return request, nil
}
func DecodeAuthUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request AuthUserRequest
	return request, nil
}
