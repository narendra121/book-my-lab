package server

import (
	"booking.com/internal/config"
	"booking.com/internal/handlers/auth"
	"booking.com/internal/handlers/properties"
	"booking.com/internal/handlers/user"
	"booking.com/internal/handlers/visits"
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
		registerNoAuthApis(v1NoAuth, cfg)
	}
	// EndPoints withAuth
	{
		v1Auth := router.Group("/v1")
		v1Auth.Use(middleware.AuthMiddleWare())

		registerUserApp(v1Auth, cfg)
		registerPropertyApp(v1Auth, cfg)
		registerVisitsApp(v1Auth, cfg)
	}

	if err := router.Run(cfg.HttpServer.Address); err != nil {
		return err
	}
	return nil
}
func registerNoAuthApis(router *gin.RouterGroup, cfg *config.AppConfig) {
	authHandler := auth.NewAuthHandler(&svcs.AuthSvc{AppCfg: cfg}, &svcs.UserSvc{AppCfg: cfg})

	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/refresh", authHandler.Refresh)
	router.POST("/auth/logout", authHandler.LogOut)
	router.PATCH("/auth/activate", authHandler.ActivateUser)

	prptyHandler := properties.NewPropertyHandler(&svcs.PropertySvc{AppCfg: cfg})
	router.GET("/properties/all", prptyHandler.GetAllProperties)
}

func registerUserApp(router *gin.RouterGroup, cfg *config.AppConfig) {
	usrHandler := user.NewUserHandler(&svcs.UserSvc{AppCfg: cfg})

	router.GET("/user/profile", usrHandler.GetProfile)
	router.PUT("/user/update", usrHandler.UpdateUser)
	router.GET("/user/list", usrHandler.ListUsers)
	router.PATCH("/user/update-role", usrHandler.UpdateRole)
	router.DELETE("/user/profile", usrHandler.DeleteUser)
}
func registerPropertyApp(router *gin.RouterGroup, cfg *config.AppConfig) {
	prptyHandler := properties.NewPropertyHandler(&svcs.PropertySvc{AppCfg: cfg})

	router.POST("/properties", prptyHandler.AddProperties)
	router.PUT("/properties", prptyHandler.UpdateProperty)
	router.GET("/properties", prptyHandler.GetFilteredProperties)
	router.DELETE("/properties/:id", prptyHandler.GetAllProperties)
}

func registerVisitsApp(router *gin.RouterGroup, cfg *config.AppConfig) {
	visitHandler := visits.NewVisitsHandler(&svcs.VisitsSvc{AppCfg: cfg})
	
	router.POST("/visits", visitHandler.ScheduleVisit)
	router.PUT("/visits", visitHandler.UpdateVisit)
	router.GET("/visits", visitHandler.FilterVisits)
	router.DELETE("/visits/:id", visitHandler.DeleteVisit)
}
