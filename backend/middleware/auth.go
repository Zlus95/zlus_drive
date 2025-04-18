package middleware

import (
	"backend/models"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

const UserContextKey = "user"

// type contextKey string

// const (
// 	UserContextKey contextKey = "user"
// )

// func RegMiddlware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var user models.User

// 		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 			http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 			return
// 		}
// 		defer r.Body.Close()

// 		if user.Email == "" || !isValidEmail(user.Email) {
// 			http.Error(w, "Invalid email", http.StatusBadRequest)
// 			return
// 		}

// 		if len(user.Password) < 5 {
// 			http.Error(w, "Password must be at least 5 characters", http.StatusBadRequest)
// 			return
// 		}

// 		if len(user.Name) < 2 {
// 			http.Error(w, "Name must be at least 2 characters", http.StatusBadRequest)
// 			return
// 		}

// 		if len(user.LastName) < 2 {
// 			http.Error(w, "Last Name must be at least 2 characters", http.StatusBadRequest)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), UserContextKey, user)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	}
// }

func RegMiddlware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if c.Request.Method == "POST" && c.Request.URL.Path == "/register" {
			if err := c.ShouldBindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Invalid JSON": err.Error()})
				c.Abort()
				return
			}

			if user.Email == "" || !isValidEmail(user.Email) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
				c.Abort()
				return
			}

			if len(user.Password) < 5 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 5 characters"})
				c.Abort()
				return
			}

			if len(user.Name) < 2 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Name must be at least 2 characters"})
				c.Abort()
				return
			}

			if len(user.LastName) < 2 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Last Name must be at least 2 characters"})
				c.Abort()
				return
			}
			c.Set(UserContextKey, user)
		}
		c.Next()
	}
}

// func LoginMiddlware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var user models.User

// 		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 			http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 			return
// 		}
// 		defer r.Body.Close()

// 		if user.Email == "" || !isValidEmail(user.Email) {
// 			http.Error(w, "Invalid email", http.StatusBadRequest)
// 			return
// 		}

// 		if len(user.Password) < 5 {
// 			http.Error(w, "Password must be at least 5 characters", http.StatusBadRequest)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), UserContextKey, user)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	}
// }
