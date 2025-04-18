package main

import (
	"backend/config"
	"backend/handlers"
	"backend/middleware"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()
	config.InitCollections(config.DB)

	r := gin.Default()

	r.Use(middleware.DevelopmentCORS())
	r.Use(middleware.RegMiddlware())
	r.Use(middleware.LoginMiddlware())

	r.POST("/register", handlers.Register)
	r.POST("login", handlers.Login)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
