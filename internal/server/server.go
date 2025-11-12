package server

import (
	"log"

	"booking.com/internal/config"
	"booking.com/internal/handlers/auth"
	"booking.com/internal/handlers/user"
	"booking.com/internal/server/middleware"
	"booking.com/internal/svcs"
	"github.com/gin-gonic/gin"
)

func StartHttpTlsServer(cfg *config.AppConfig) error {
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
	log.Printf("Server started on https://%v", cfg.HttpServer.Address)
	if err := router.Run(cfg.HttpServer.Address); err != nil {
		return err
	}
	return nil
}

func registerUsersAppWithAuth(router *gin.RouterGroup, cfg *config.AppConfig) {
	usrHandler := user.NewUserHandler(&svcs.UserSvc{AppCfg: cfg})

	router.POST("/user/register", usrHandler.Add)
	router.GET("/user/profile/:username", usrHandler.Get)
	router.PUT("/user/update", usrHandler.Put)
	router.GET("/user/list", usrHandler.GetAll)
	router.PUT("/user/update-role", usrHandler.UpdateRole)
}

func registerAuthAppNoAuth(router *gin.RouterGroup, cfg *config.AppConfig) {
	authHandler := auth.NewAuthHandler(&svcs.AuthSvc{AppCfg: cfg})

	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/refresh", authHandler.Refresh)
	router.POST("/auth/logout", authHandler.LogOut)
}
