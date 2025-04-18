package handlers

import (
	"backend/config"
	"backend/middleware"
	"backend/models"
	"backend/utils"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userValue, ok := c.Get(middleware.UserContextKey)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}

	user, ok := userValue.(models.User)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data type"})
		return
	}

	userCollection := config.UserCollection

	foundUser := bson.M{"email": user.Email}

	err := userCollection.FindOne(ctx, foundUser).Err()
	if err == nil {
		log.Printf("User with email %s already exists", user.Email)
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		log.Printf("Error checking user existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence"})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	user.Password = string(hashPassword)
	user.StorageLimit = 1 * 1024 * 1024 * 1024
	user.UsedStorage = 0

	if _, err := userCollection.InsertOne(ctx, user); err != nil {
		log.Printf("Error inserting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           user.ID.Hex(),
		"name":         user.Name,
		"lastName":     user.LastName,
		"storageLimit": user.StorageLimit,
		"usedStorage":  user.UsedStorage,
	})
}

func Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userValue, ok := c.Get(middleware.UserContextKey)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}

	user, ok := userValue.(models.User)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data type"})
		return
	}

	userCollection := config.UserCollection

	var currentUser models.User

	if err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&currentUser); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("User not found: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect Email"})
			return
		} else if err != mongo.ErrNoDocuments {
			log.Printf("Error finding user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(user.Password)); err != nil {
		log.Printf("Invalid password for user %s: %v", user.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect Password"})
		return
	}

	token, err := utils.GenerateJWT(currentUser.ID.Hex())

	if err != nil {
		log.Printf("JWT generation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}
