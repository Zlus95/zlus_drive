package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	r.HandleFunc("/register", middleware.RegMiddlware(handlers.Register)).Methods("POST")
	r.HandleFunc("/login", middleware.LoginMiddlware(handlers.Login)).Methods("POST")
}
