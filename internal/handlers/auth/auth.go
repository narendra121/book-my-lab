package auth

import (
	"errors"
	"net/http"

	"booking.com/internal/dto"
	"booking.com/internal/svcs"
	"booking.com/internal/utils"
	jwtauth "booking.com/pkg/auth/jwt-auth"
	"booking.com/pkg/constants"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	AuthSvc *svcs.AuthSvc
	UsrSvc  *svcs.UserSvc
}

func NewAuthHandler(authSvc *svcs.AuthSvc, usrSvc *svcs.UserSvc) *AuthHandler {
	return &AuthHandler{
		AuthSvc: authSvc,
		UsrSvc:  usrSvc,
	}
}
func (a *AuthHandler) Register(c *gin.Context) {
	var userReq dto.CreateUser
	if err := c.ShouldBindBodyWithJSON(&userReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", err, nil))
		return
	}
	var errMsg string
	if userReq.Phone == "" {
		errMsg += "phone number not provided"
	}
	if errMsg != "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", errors.New(errMsg), nil))
		return
	}
	if err := a.AuthSvc.RegisterUser(&userReq, a.UsrSvc); err != nil {
		if errors.Is(err, utils.ErrUserAlreadyExistsWithEmail) || errors.Is(err, utils.ErrUserAlreadyExistsWithPhone) {
			c.AbortWithStatusJSON(http.StatusOK, utils.WriteAppResponse("", err, nil))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	c.JSON(http.StatusCreated, utils.WriteAppResponse("user registered successfully.", nil, nil))
}

func (a *AuthHandler) Login(c *gin.Context) {
	var reqUser dto.Login
	if err := c.ShouldBindBodyWithJSON(&reqUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", err, nil))
		return
	}
	if reqUser.UserName == "" || reqUser.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", errors.New("username or password is empty"), nil))
		return
	}
	token, refreshToken, err := a.AuthSvc.Login(reqUser, a.UsrSvc)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.WriteAppResponse("", utils.ErrUserNotFound, nil))
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.WriteAppResponse("", err, nil))
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     constants.RefreshToken,
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})
	c.JSON(http.StatusOK, utils.WriteAppResponse("", nil, map[string]string{
		constants.AccessToken: token,
		constants.TokenType:   constants.Bearer,
	}))
}

func (a *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(constants.RefreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", err, nil))
		return
	}
	if refreshToken == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", errors.New("missing refresh token"), nil))
		return
	}
	userName, err := jwtauth.GetUnVerifiedJwtClaims(refreshToken, constants.UserName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, utils.WriteAppResponse("", err, nil))
		return
	}
	if userName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("invalid refresh token", nil, nil))
		return
	}

	newToken, newRefreshToken, err := a.AuthSvc.Refresh(userName, refreshToken, a.UsrSvc)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.WriteAppResponse("", utils.ErrUserNotFound, nil))
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.WriteAppResponse("", err, nil))
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     constants.RefreshToken,
		Value:    newRefreshToken,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})
	c.JSON(http.StatusOK, utils.WriteAppResponse("", nil, map[string]string{
		constants.AccessToken: newToken,
		constants.TokenType:   constants.Bearer,
	}))
}
func (a *AuthHandler) LogOut(c *gin.Context) {
	refreshToken, err := c.Cookie(constants.RefreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", err, nil))
		return
	}
	if refreshToken == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", errors.New("missing refresh_token"), nil))
		return
	}

	userName, err := jwtauth.GetUnVerifiedJwtClaims(refreshToken, constants.UserName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", err, nil))
		return
	}
	if userName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", errors.New("invalid token data"), nil))
		return
	}
	if err := a.AuthSvc.LogOut(userName, a.UsrSvc); err != nil {
		if errors.Is(err, utils.ErrUserAlreadyLoggedOut) {
			c.AbortWithStatusJSON(http.StatusOK, utils.WriteAppResponse("", utils.ErrUserAlreadyLoggedOut, nil))
			return
		}
		if errors.Is(err, utils.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.WriteAppResponse("", utils.ErrUserNotFound, nil))
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.WriteAppResponse("", utils.ErrUserNotFound, nil))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	c.JSON(http.StatusOK, utils.WriteAppResponse("user logged out successfully", nil, nil))
}
func (a *AuthHandler) ActivateUser(c *gin.Context) {
	var userReq dto.Activate
	if err := c.ShouldBindBodyWithJSON(&userReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", err, nil))
		return
	}
	if err := a.AuthSvc.ActivateUser(userReq.UserName, a.UsrSvc); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	c.JSON(http.StatusOK, utils.WriteAppResponse("user activated successfully", nil, nil))
}
