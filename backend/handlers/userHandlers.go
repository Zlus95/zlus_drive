package handlers

import (
	"backend/config"
	"backend/middleware"
	"backend/models"
	"backend/utils"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GetCurrentUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userIDValue, ok := c.Get(middleware.UserIDKey)

	if !ok || userIDValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}

	userID, ok := userIDValue.(string)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	userCollection := config.UserCollection
	var user models.User

	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)

	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Database request timed out"})

		} else if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	response := map[string]interface{}{
		"id":           user.ID.Hex(),
		"name":         user.Name,
		"lastName":     user.LastName,
		"email":        user.Email,
		"storageLimit": user.StorageLimit,
		"usedStorage":  user.UsedStorage,
	}

	c.JSON(http.StatusOK, response)
}

func ChangeStorageLimit(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userIDValue, ok := c.Get(middleware.UserIDKey)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be a string"})
		return
	}

	userIDHex, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format: " + err.Error()})
		return
	}

	userLimitValue, ok := c.Get(middleware.StorageLimitContextKey)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Storage limit not found in context"})
		return
	}

	userLimit, ok := userLimitValue.(int64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid storage limit type, expected int64"})
		return
	}

	if userLimit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Storage limit cannot be negative"})
		return
	}

	userCollection := config.UserCollection

	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": userIDHex}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			log.Printf("Failed to find user %s: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user: " + err.Error()})
		}
		return
	}

	update := bson.M{"$set": bson.M{"storageLimit": userLimit}}
	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": userIDHex}, update)
	if err != nil {
		log.Printf("Failed to update storage limit for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update storage limit: " + err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found for update"})
		return
	}

	if result.ModifiedCount == 0 {
		log.Printf("Storage limit for user %s was not modified, possibly same value", userID)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Storage limit changed successfully",
		"newLimit": userLimit,
	})
}
