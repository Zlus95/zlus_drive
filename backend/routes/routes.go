package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(middleware.DevelopmentCORS())
	r.Use(middleware.RegMiddlware())
	r.Use(middleware.LoginMiddlware())
	r.Use(middleware.TokenMiddlware())

	r.POST("/register", handlers.Register)
	r.POST("login", handlers.Login)
}
