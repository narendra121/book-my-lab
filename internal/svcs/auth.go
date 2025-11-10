package svcs

import (
	"errors"
	"fmt"

	"booking.com/internal/config"
	"booking.com/internal/db/postgresql/dao"
	"booking.com/internal/dto"
	jwtauth "booking.com/pkg/auth/jwt-auth"
	"booking.com/pkg/utils"
)

type AuthSvc struct {
	AppCfg *config.AppConfig
}

func (a *AuthSvc) Login(reqUser dto.Login) (string, error) {
	q := dao.User
	users, err := q.Where(q.Email.Eq(reqUser.UserName)).Or(q.Phone.Eq(reqUser.UserName)).Find()
	if err != nil || len(users) == 0 {
		return "", errors.New("user not found")
	}
	user := users[0]
	if validPassword := utils.CheckPassword(user.PasswordHash, reqUser.Password+user.Salt); !validPassword {
		return "", errors.New("invalid password")
	}
	token, err := jwtauth.GetToken(user.Email, user.Salt, a.AppCfg.Jwt.Expiry)
	if err != nil {
		return "", fmt.Errorf("failed to generate token, error: %v", err)
	}
	return token, nil
}

func (a *AuthSvc) Refresh(tokenReq dto.TokenReq) (string, error) {
	newToken, err := jwtauth.RefreshToken(tokenReq.AccessToken, a.AppCfg.Jwt.Expiry, true)
	if err != nil {
		return "", fmt.Errorf("token refresh failed: %v", err)
	}
	return newToken, nil
}
