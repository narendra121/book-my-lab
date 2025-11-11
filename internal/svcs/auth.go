package svcs

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"booking.com/internal/config"
	"booking.com/internal/db/postgresql/dao"
	"booking.com/internal/db/postgresql/model"
	"booking.com/internal/dto"
	jwtauth "booking.com/pkg/auth/jwt-auth"
	"booking.com/pkg/constants"
	"booking.com/pkg/utils"
)

type AuthSvc struct {
	AppCfg *config.AppConfig
}

func (a *AuthSvc) Login(reqUser dto.Login) (string, string, error) {
	q := dao.User
	user, err := q.Where(q.Email.Eq(reqUser.UserName)).Or(q.Phone.Eq(reqUser.UserName)).First()
	if err != nil || user == nil {
		return "", "", errors.New("user not found")
	}
	if validPassword := utils.CheckPassword(user.PasswordHash, reqUser.Password+user.Salt); !validPassword {
		return "", "", errors.New("invalid password")
	}
	return a.getAccessAndRefreshTokens(user)
}

func (a *AuthSvc) Refresh(refreshToken string) (string, string, error) {
	userName, err := jwtauth.GetUnVerifiedJwtClaims(refreshToken, constants.UserName)
	if err != nil || userName == "" {
		return "", "", fmt.Errorf("invalid token data")
	}
	q := dao.User
	user, _ := q.Where(q.Email.Eq(userName)).Or(q.Phone.Eq(userName)).First()
	if user == nil {
		return "", "", fmt.Errorf("invalid token data")
	}

	validateFunc := func(userName ...string) bool {
		q := dao.User
		user, _ := q.Where(q.Email.Eq(userName[0])).Or(q.Phone.Eq(userName[0])).First()
		return user != nil
	}

	validRefreshToken, err := jwtauth.IsTokenValid(refreshToken, user.Salt, validateFunc)
	if err != nil || !validRefreshToken {
		return "", "", fmt.Errorf("invalid refresh_token")
	}

	hash := sha256.Sum256([]byte(refreshToken))
	if user == nil || user.RefreshToken != hex.EncodeToString(hash[:]) {
		return "", "", fmt.Errorf("refresh_token revoked")
	}
	return a.getAccessAndRefreshTokens(user)
}

func (a *AuthSvc) getAccessAndRefreshTokens(user *model.User) (string, string, error) {
	token, err := jwtauth.GetToken(user.Email, user.Salt, a.AppCfg.Jwt.AccessTokenExpiry)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate token, error: %v", err)
	}
	refreshToken, err := jwtauth.GetToken(user.Email, user.Salt, a.AppCfg.Jwt.RefreshTokenExpiry)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate token, error: %v", err)
	}
	hash := sha256.Sum256([]byte(refreshToken))
	q := dao.User
	_, err = q.Where(q.Email.Eq(user.Email)).Updates(&model.User{RefreshToken: hex.EncodeToString(hash[:])})
	if err != nil {
		return "", "", fmt.Errorf("failed to store refresh_token, error: %v", err)
	}
	return token, refreshToken, nil
}
