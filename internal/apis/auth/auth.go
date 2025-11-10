package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"booking.com/internal/db/postgresql/dao"
	"booking.com/internal/dto"
	jwtauth "booking.com/pkg/auth/jwt-auth"
	"booking.com/pkg/constants"
	"booking.com/pkg/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
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
	q := dao.User
	users, err := q.Where(q.Email.Eq(reqUser.UserName)).Or(q.Phone.Eq(reqUser.UserName)).Find()
	if err != nil || len(users) == 0 {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	user := users[0]
	if validPassword := utils.CheckPassword(user.PasswordHash, reqUser.Password+user.Salt); !validPassword {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}
	token, err := jwtauth.GetToken(user.Email, user.Salt, 15)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to generate token, error: %v", err), http.StatusUnauthorized)
		return
	}
	w.Header().Set(constants.ContentType, constants.ContentTypeJson)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		constants.AccessToken: token,
		constants.TokenType:   constants.Bearer,
	})
}

func Refresh(w http.ResponseWriter, r *http.Request) {
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
	newToken, err := jwtauth.RefreshToken(tokenReq.AccessToken, true)
	if err != nil {
		http.Error(w, fmt.Sprintf("token refresh failed: %v", err), http.StatusUnauthorized)
		return
	}
	w.Header().Set(constants.ContentType, constants.ContentTypeJson)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		constants.AccessToken: newToken,
		constants.TokenType:   constants.Bearer,
	})
}
