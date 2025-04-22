package handlers

import (
	"backend/config"
	"backend/middleware"
	"backend/models"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddFile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userIDValue, ok := c.Get(middleware.UserIDKey)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in context"})
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

	fileInfoValue, ok := c.Get(middleware.FileContextKey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File info not found in context"})
		return
	}

	fileInfo, ok := fileInfoValue.(middleware.FileInfo)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid file info type"})
		return
	}

	fileSizeMB := int((fileInfo.Size + (1<<20 - 1)) / (1 << 20))
	if fileSizeMB < 1 {
		fileSizeMB = 1
	}

	var user models.User
	if err := config.UserCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data: " + err.Error()})
		return
	}

	if user.UsedStorage+fileSizeMB > user.StorageLimit {
		c.JSON(http.StatusForbidden, gin.H{
			"error":        "Storage limit exceeded",
			"usedStorage":  user.UsedStorage,
			"storageLimit": user.StorageLimit,
			"fileSize":     fileSizeMB,
		})
		return
	}

	fileID := primitive.NewObjectID()
	filePath := fmt.Sprintf("uploads/%s/%s", userID, fileID.Hex())
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory: " + err.Error()})
		return
	}

	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file: " + err.Error()})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, fileInfo.File); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file: " + err.Error()})
		return
	}

	updateResult, err := config.UserCollection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"usedStorage": fileSizeMB}},
	)
	if err != nil {
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user storage: " + err.Error()})
		return
	}

	if updateResult.ModifiedCount == 0 {
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User storage not updated"})
		return
	}

	newFile := models.File{
		ID:        fileID,
		Name:      fileInfo.Filename,
		Size:      fileInfo.Size,
		MimeType:  fileInfo.MimeType,
		Path:      filePath,
		OwnerID:   objID,
		CreatedAt: time.Now(),
	}

	if _, err := config.FilesCollection.InsertOne(ctx, newFile); err != nil {
		config.UserCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": objID},
			bson.M{"$inc": bson.M{"usedStorage": -fileSizeMB}},
		)
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata: " + err.Error()})
		return
	}

	var updatedUser models.User
	if err := config.UserCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedUser); err == nil {
		user.UsedStorage = updatedUser.UsedStorage
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "File uploaded successfully",
		"fileId":       fileID.Hex(),
		"name":         fileInfo.Filename,
		"size":         fileInfo.Size,
		"type":         fileInfo.MimeType,
		"usedStorage":  user.UsedStorage,
		"storageLimit": user.StorageLimit,
	})
}
