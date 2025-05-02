package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func FileRoutes(r *gin.Engine) {
	r.Use(middleware.TokenMiddlware())

	r.POST("/upload",
		middleware.FileMiddlware,
		handlers.AddFile,
	)
	r.DELETE("/file/:id", handlers.DeleteFile)

	r.GET("/files", handlers.GetAllFiles)

	r.POST("/folder",
		middleware.FolderMiddlware(),
		handlers.CreateFolder,
	)

	r.PATCH("/file/:id", handlers.MoveFile)

}
