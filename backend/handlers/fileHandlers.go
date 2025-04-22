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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// добавить проверку для StorageLimit && UsedStorage
func AddFile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)

	defer cancel()

	userIDValue, ok := c.Get(middleware.UserIDKey)

	if !ok || userIDValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}

	userID, ok := userIDValue.(string)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data type"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	file.Seek(0, 0)

	mimeType := http.DetectContentType(buffer)

	fileID := primitive.NewObjectID()

	filePath := fmt.Sprintf("uploads/%s/%s", userID, fileID.Hex())

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	newFile := models.File{
		ID:        fileID,
		Name:      header.Filename,
		Size:      header.Size,
		MimeType:  mimeType,
		Path:      filePath,
		OwnerID:   objID,
		CreatedAt: time.Now(),
	}

	if _, err := config.FilesCollection.InsertOne(ctx, newFile); err != nil {
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"fileId":  fileID.Hex(),
		"name":    header.Filename,
		"size":    header.Size,
		"type":    mimeType,
	})
}
