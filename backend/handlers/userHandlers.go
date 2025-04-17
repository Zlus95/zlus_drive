package handlers

import (
	"backend/config"
	"backend/middleware"
	"backend/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	user, ok := r.Context().Value(middleware.UserContextKey).(models.User)

	if !ok {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	userCollection := config.UserCollection

	foundUser := bson.M{"email": user.Email}

	err := userCollection.FindOne(ctx, foundUser).Err()
	if err == nil {
		log.Printf("User with email %s already exists", user.Email)
		http.Error(w, "User already exists", http.StatusConflict)
		return
	} else if err != mongo.ErrNoDocuments {
		log.Printf("Error checking user existence: %v", err)
		http.Error(w, "Failed to check user existence", http.StatusInternalServerError)
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashPassword)
	user.StorageLimit = 1 * 1024 * 1024 * 1024
	user.UsedStorage = 0

	if _, err := userCollection.InsertOne(ctx, user); err != nil {
		log.Printf("Error inserting user: %v", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":           user.ID.Hex(),
		"name":         user.Name,
		"lastName":     user.LastName,
		"storageLimit": user.StorageLimit,
		"usedStorage":  user.UsedStorage,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}
