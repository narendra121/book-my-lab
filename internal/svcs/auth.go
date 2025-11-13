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
	usr, err := userSvc.GetUserWithEmailOrPhone(userReq.Email, userReq.Phone, false)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if usr != nil {
		var usrE []*model.User
		if userReq.Email != "" {
			usrE, err = userSvc.FilterUsers("", userReq.Email, "", false)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
		}
		usrP, err := userSvc.FilterUsers("", "", userReq.Phone, false)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if len(usrE) > 0 && len(usrP) > 0 {
			if usrE[0].Deleted && usrP[0].Deleted {
				return errors.New(utils.ErrUserAlreadyExistsWithEmailAndPhone.Error() + ". Please activate")
			}
			return utils.ErrUserAlreadyExistsWithEmailAndPhone
		}
		if len(usrE) > 0 && len(usrP) == 0 {
			if usrE[0].Deleted {
				return errors.New(utils.ErrUserAlreadyExistsWithEmail.Error() + ". Please activate")
			}
			return utils.ErrUserAlreadyExistsWithEmail
		} else {
			if usrP[0].Deleted {
				return errors.New(utils.ErrUserAlreadyExistsWithPhone.Error() + ". Please activate")
			}
			return utils.ErrUserAlreadyExistsWithPhone
		}
	}
	salt := utils.GetUUID()
	hashedPass, err := utils.HashPassword(userReq.Password + salt)
	if err != nil {
		return err
	}
	err = userSvc.CreateUser(&model.User{
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
	})
	if err != nil {
		return err
	}
	return nil
}
func (a *AuthSvc) Login(reqUser dto.Login, userSvc *UserSvc) (string, string, error) {
	user, err := userSvc.GetUserWithEmailOrPhone(reqUser.UserName, reqUser.UserName, true)
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
		user, _ := q.Where(q.Username.Eq(userName[0])).Where(q.Deleted.Is(false)).First()
		return user != nil
	}
	user, err := userSvc.GetUserByUserName(userName, true)
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
	user, err := usrSvc.GetUserByUserName(userName, true)
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
	token, err := jwtauth.GetToken(user.Username, user.Salt, a.AppCfg.Jwt.AccessTokenExpiry)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwtauth.GetToken(user.Username, user.Salt, a.AppCfg.Jwt.RefreshTokenExpiry)
	if err != nil {
		return "", "", err
	}
	hash := sha256.Sum256([]byte(refreshToken))
	err = userSvc.UpdateRefreshToken(user.Username, hex.EncodeToString(hash[:]))
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}
func (a *AuthSvc) ActivateUser(userName string, usrSvc *UserSvc) error {
	user, err := usrSvc.GetUserWithEmailOrPhone(userName, userName, true)
	if err != nil {
		return err
	}
	if user == nil {
		return utils.ErrUserAlreadyExistsWithEmailOrPhone
	}
	return usrSvc.UpdateDelFlag(user.Username, false)
}
