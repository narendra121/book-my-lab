package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"booking.com/internal/svcs"
	"booking.com/internal/utils"
	jwtauth "booking.com/pkg/auth/jwt-auth"
	"booking.com/pkg/constants"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, utils.WriteAppResponse("Sever is Up and Running", nil, nil))
}
func CommonChain() gin.HandlersChain {
	return []gin.HandlerFunc{
		gin.Recovery(),
		corsMiddleware(),
		logFormatMiddleWare(),
	}
}

func corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:8080",
		},
		AllowCredentials: true,
	})
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get(constants.Authorization)
		if authHeader == "" || !strings.HasPrefix(authHeader, constants.Bearer) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.WriteAppResponse("", errors.New("bearer token not provided"), nil))
			return
		}
		token := authHeader[len(constants.Bearer):]
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.WriteAppResponse("", errors.New("access token not provided"), nil))
			return
		}
		userName, err := jwtauth.GetUnVerifiedJwtClaims(token, constants.UserName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.WriteAppResponse("", err, nil))
			return
		}
		if userName == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.WriteAppResponse("", errors.New("invalid token claims"), nil))
			return
		}
		usrSvc := &svcs.UserSvc{}
		user, err := usrSvc.GetUserWithDelFlag(userName)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.WriteAppResponse("", utils.ErrUserNotFound, nil))
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.WriteAppResponse("", err, nil))
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.WriteAppResponse("", utils.ErrUserNotFound, nil))
			return
		}
		validToken, err := jwtauth.IsTokenValid(token, user.Salt, nil)
		if err != nil || !validToken || user.RefreshToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.WriteAppResponse("", errors.New("session expired"), nil))
			return
		}
		c.Set(constants.CurrentUser, user)
		c.Set(constants.Role, user.Role)
		c.Set(constants.CurrentUserName, userName)
		c.Next()
	}
}
func logFormatMiddleWare() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}
