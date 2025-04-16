package main

import (
	"backend/config"
	"backend/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()
	config.InitCollections(config.DB)

	r := mux.NewRouter()

	routes.UserRoutes(r)

	// Настройка CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Запуск сервера
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", c.Handler(r)))
}
