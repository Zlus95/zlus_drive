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

func LoginMiddlware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if c.Request.Method == "POST" && c.Request.URL.Path == "/login" {
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
			c.Set(UserContextKey, user)
		}
		c.Next()
	}
}
