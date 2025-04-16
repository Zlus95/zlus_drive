package middleware

import (
	"backend/models"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

func RegMiddlware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if user.Email == "" || !strings.Contains(user.Email, "@") {
			http.Error(w, "Invalid email", http.StatusBadRequest)
			return
		}

		if len(user.Password) < 5 {
			http.Error(w, "Password must be at least 5 characters", http.StatusBadRequest)
			return
		}

		if len(user.Name) < 2 {
			http.Error(w, "Name must be at least 2 characters", http.StatusBadRequest)
			return
		}

		if len(user.LastName) < 2 {
			http.Error(w, "Last Name must be at least 2 characters", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
