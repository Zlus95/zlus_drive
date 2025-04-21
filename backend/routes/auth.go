package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/register",
		middleware.RegMiddlware(),
		handlers.Register,
	)

	r.POST("/login",
		middleware.LoginMiddlware(),
		handlers.Login,
	)

	r.GET("/user",
		middleware.TokenMiddlware(),
		handlers.GetCurrentUser,
	)

}
