package auth

import (
	"errors"
	"net/http"

	"booking.com/internal/dto"
	"booking.com/internal/svcs"
	"booking.com/internal/utils"
	"booking.com/pkg/constants"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthSvc *svcs.AuthSvc
}

func NewAuthHandler(authSvc *svcs.AuthSvc) *AuthHandler {
	return &AuthHandler{
		AuthSvc: authSvc,
	}
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
	token, refreshToken, err := a.AuthSvc.Login(reqUser)
	if err != nil {
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
	newToken, newRefreshToken, err := a.AuthSvc.Refresh(refreshToken)
	if err != nil {
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
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.WriteAppResponse("", errors.New("missing refresh token"), nil))
		return
	}
	err = a.AuthSvc.LogOut(refreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.WriteAppResponse("", err, nil))
		return
	}
	c.JSON(http.StatusOK, utils.WriteAppResponse("user logged out successfully", nil, nil))
}
