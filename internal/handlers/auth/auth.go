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
	token, refreshToken, err := a.AuthSvc.Login(reqUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RefreshToken,
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})
	w.Header().Set(constants.ContentType, constants.ContentTypeJson)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		constants.AccessToken: token,
		constants.TokenType:   constants.Bearer,
	})
}

func (a *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(constants.RefreshToken)
	if err != nil {
		http.Error(w, "missing refresh token", http.StatusUnauthorized)
		return
	}
	refreshToken := cookie.Value
	if refreshToken == "" {
		http.Error(w, "access token is missing", http.StatusBadRequest)
		return
	}
	newToken, newRefreshToken, err := a.AuthSvc.Refresh(refreshToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RefreshToken,
		Value:    newRefreshToken,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})
	w.Header().Set(constants.ContentType, constants.ContentTypeJson)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		constants.AccessToken: newToken,
		constants.TokenType:   constants.Bearer,
	})
}
