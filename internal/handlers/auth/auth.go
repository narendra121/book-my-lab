package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"booking.com/internal/dto"
	"booking.com/internal/svcs"
	"booking.com/pkg/constants"
	"booking.com/pkg/utils"
)

type AuthHandler struct {
	AuthSvc *svcs.AuthSvc
}

func NewAuthHandler(authSvc *svcs.AuthSvc) *AuthHandler {
	return &AuthHandler{
		AuthSvc: authSvc,
	}
}
func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var reqUser dto.Login
	err := utils.ParseHttpRequest(r, &reqUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid reqBody, error: %v", err), http.StatusBadRequest)
		return
	}
	if reqUser.UserName == "" || reqUser.Password == "" {
		http.Error(w, "username/password is empty", http.StatusBadRequest)
		return
	}
	token, err := a.AuthSvc.Login(reqUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusUnauthorized)
		return
	}
	w.Header().Set(constants.ContentType, constants.ContentTypeJson)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		constants.AccessToken: token,
		constants.TokenType:   constants.Bearer,
	})
}

func (a *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var tokenReq dto.TokenReq
	err := utils.ParseHttpRequest(r, &tokenReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid reqBody, error: %v", err), http.StatusBadRequest)
		return
	}
	if tokenReq.AccessToken == "" {
		http.Error(w, "access token is missing", http.StatusBadRequest)
		return
	}
	newToken, err := a.AuthSvc.Refresh(tokenReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusUnauthorized)
		return
	}
	w.Header().Set(constants.ContentType, constants.ContentTypeJson)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		constants.AccessToken: newToken,
		constants.TokenType:   constants.Bearer,
	})
}
