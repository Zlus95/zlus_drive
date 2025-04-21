package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.Use(middleware.TokenMiddlware())

	r.GET("/user", handlers.GetCurrentUser)
}
