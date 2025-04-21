package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.Use(middleware.DevelopmentCORS())
	r.Use(middleware.RegMiddlware())
	r.Use(middleware.LoginMiddlware())

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
}
