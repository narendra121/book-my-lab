package server

import (
	"booking.com/internal/config"
	"booking.com/internal/handlers/auth"
	"booking.com/internal/handlers/user"
	"booking.com/internal/server/middleware"
	"booking.com/internal/svcs"
	"github.com/gin-gonic/gin"
)

func StartHttpTlsServer(cfg *config.AppConfig) error {
	gin.SetMode(cfg.HttpServer.Mode)

	router := gin.New()

	router.Use(middleware.CommonChain()...)

	router.GET("/health", middleware.Health)
	// EndPoints withoutAuth
	{
		v1NoAuth := router.Group("/v1")
		registerAuthAppNoAuth(v1NoAuth, cfg)
	}
	// EndPoints withAuth
	{
		v1Auth := router.Group("/v1")
		v1Auth.Use(middleware.AuthMiddleWare())

		registerUsersAppWithAuth(v1Auth, cfg)
	}

	if err := router.Run(cfg.HttpServer.Address); err != nil {
		return err
	}
	return nil
}

func registerUsersAppWithAuth(router *gin.RouterGroup, cfg *config.AppConfig) {
	usrHandler := user.NewUserHandler(&svcs.UserSvc{AppCfg: cfg})

	router.GET("/user/profile", usrHandler.GetProfile)
	router.PUT("/user/update", usrHandler.UpdateUser)
	router.GET("/user/list", usrHandler.ListUsers)
	router.PUT("/user/update-role", usrHandler.UpdateRole)
	router.DELETE("/user/profile", usrHandler.DeleteUser)
}

func registerAuthAppNoAuth(router *gin.RouterGroup, cfg *config.AppConfig) {
	authHandler := auth.NewAuthHandler(&svcs.AuthSvc{AppCfg: cfg}, &svcs.UserSvc{AppCfg: cfg})

	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/refresh", authHandler.Refresh)
	router.POST("/auth/logout", authHandler.LogOut)
}
