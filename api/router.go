package api

import (
	"auth_service/api/handlers"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
	_ "auth_service/api/docs"
)

// Router ...
// @title           Task Management System Auth
// @version         1
// @description     Task Management System Auth
// @in header
// @name Authorization
func Router(h *handlers.Handler) *gin.Engine {
	router := gin.New()

	router.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		RequestHeaders: "Authorization,Origins,Content-Type",
		Methods:        "POST, GET, PUT, DELETE, OPTIONS",
	}))

	auth_service := router.Group("/auth_service")
	{
		auth_service.POST("/register", h.Reginster)
		auth_service.POST("/login", h.Login)
		auth_service.POST("/reset_password", h.ResetPassword)
		auth_service.POST("/change_password", h.ChangePassword)
		auth_service.POST("/forgot_password", h.ForgotPassword)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
