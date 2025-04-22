package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func FileRoutes(r *gin.Engine) {
	r.Use(middleware.TokenMiddlware())

	r.POST("/upload",
		handlers.AddFile,
	)

}
