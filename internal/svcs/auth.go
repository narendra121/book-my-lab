package svcs

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"booking.com/internal/config"
	"booking.com/internal/db/postgresql/dao"
	"booking.com/internal/db/postgresql/model"
	"booking.com/internal/dto"
	"booking.com/internal/utils"
	jwtauth "booking.com/pkg/auth/jwt-auth"
	"booking.com/pkg/constants"
	"gorm.io/gorm"
)

type AuthSvc struct {
	AppCfg *config.AppConfig
}

func NewAuthSvc(cfg *config.AppConfig) *AuthSvc {
	return &AuthSvc{AppCfg: cfg}
}
func (a *AuthSvc) RegisterUser(userReq *dto.CreateUser, userSvc *UserSvc) error {
	user, err := userSvc.GetUserWithoutDelFlag(userReq.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if user != nil && !user.Deleted {
		return utils.ErrUserAlreadyExists
	}
	if user != nil && user.Deleted {
		err := userSvc.UpdateDelFlag(userReq.Email, false)
		if err != nil {
			return err
		}
	}
	salt := utils.GetUUID()
	hashedPass, err := utils.HashPassword(userReq.Password + salt)
	if err != nil {
		return err
	}
	err = userSvc.CreateUser([]*model.User{{
		FirstName:    userReq.FirstName,
		LastName:     userReq.LastName,
		Email:        userReq.Email,
		Phone:        userReq.Phone,
		PasswordHash: hashedPass,
		Salt:         salt,
		Address:      userReq.Address,
		Role:         constants.UserRole,
		Deleted:      false,
		UpdatedAt:    time.Now(),
	}}...)
	if err != nil {
		return err
	}
	return nil
}
func (a *AuthSvc) Login(reqUser dto.Login, userSvc *UserSvc) (string, string, error) {
	user, err := userSvc.GetUserWithDelFlag(reqUser.UserName)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", utils.ErrUserNotFound
	}
	if validPassword := utils.CheckPassword(user.PasswordHash, reqUser.Password+user.Salt); !validPassword {
		return "", "", utils.ErrInvalidUserOrPass
	}
	return a.getAccessAndRefreshTokens(user, userSvc)
}

func (a *AuthSvc) Refresh(userName, refreshToken string, userSvc *UserSvc) (string, string, error) {
	validateFunc := func(userName ...string) bool {
		q := dao.User
		user, _ := q.Where(q.Email.Eq(userName[0])).Or(q.Phone.Eq(userName[0])).Where(q.Deleted.Is(false)).First()
		return user != nil
	}
	user, err := userSvc.GetUserWithDelFlag(userName)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", utils.ErrUserNotFound
	}
	validRefreshToken, err := jwtauth.IsTokenValid(refreshToken, user.Salt, validateFunc)
	if err != nil {
		return "", "", err
	}
	if !validRefreshToken {
		return "", "", errors.New("invalid refresh_token")
	}
	hash := sha256.Sum256([]byte(refreshToken))
	if user.RefreshToken != hex.EncodeToString(hash[:]) {
		return "", "", fmt.Errorf("refresh_token revoked")
	}
	return a.getAccessAndRefreshTokens(user, userSvc)
}
func (a *AuthSvc) LogOut(userName string, usrSvc *UserSvc) error {
	user, err := usrSvc.GetUserWithDelFlag(userName)
	if err != nil {
		return err
	}
	if user == nil {
		return utils.ErrUserNotFound
	}
	if user.RefreshToken == "" {
		return utils.ErrUserAlreadyLoggedOut
	}
	err = usrSvc.UpdateRefreshToken(userName, "")
	if err != nil {
		return err
	}
	return nil
}
func (a *AuthSvc) getAccessAndRefreshTokens(user *model.User, userSvc *UserSvc) (string, string, error) {
	token, err := jwtauth.GetToken(user.Email, user.Salt, a.AppCfg.Jwt.AccessTokenExpiry)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwtauth.GetToken(user.Email, user.Salt, a.AppCfg.Jwt.RefreshTokenExpiry)
	if err != nil {
		return "", "", err
	}
	hash := sha256.Sum256([]byte(refreshToken))
	err = userSvc.UpdateRefreshToken(user.Email, hex.EncodeToString(hash[:]))
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}
