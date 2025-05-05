package handlers

import (
	"backend/config"
	"backend/middleware"
	"backend/models"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	fileSizeMB := int64((fileInfo.Size + (1<<20 - 1)) / (1 << 20))
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

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
	})
}

func DeleteFile(c *gin.Context) {
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

	fileID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID format"})
		return
	}

	var file models.File
	if err := config.FilesCollection.FindOne(
		ctx,
		bson.M{"_id": fileID, "ownerId": objID},
	).Decode(&file); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file: " + err.Error()})
		}
		return
	}

	fileSizeMB := int64((file.Size + (1<<20 - 1)) / (1 << 20))
	if fileSizeMB < 1 {
		fileSizeMB = 1
	}

	if _, err := config.FilesCollection.DeleteOne(ctx, bson.M{"_id": fileID, "ownerId": objID}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file: " + err.Error()})
		return
	}

	if !file.IsFolder {
		if _, err := config.UserCollection.UpdateOne(
			ctx,
			bson.M{"_id": objID},
			bson.M{"$inc": bson.M{"usedStorage": -fileSizeMB}},
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update storage: " + err.Error()})
			return
		}

		filePath := fmt.Sprintf("uploads/%s/%s", userID, fileID.Hex())
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			log.Printf("Failed to delete file from disk: %v", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
	})
}

func GetAllFiles(c *gin.Context) {
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

	sortByFolders := c.DefaultQuery("sort", "true") == "true"
	sortDirection := -1

	if !sortByFolders {
		sortDirection = 1
	}

	opts := options.Find().SetSort(bson.D{
		{Key: "isFolder", Value: sortDirection},
		{Key: "name", Value: 1},
	})
	cursor, err := config.FilesCollection.Find(ctx, bson.M{"ownerId": objID}, opts)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch files"})
		return
	}
	defer cursor.Close(ctx)

	var files []models.File

	if err = cursor.All(ctx, &files); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode files"})
		return
	}

	response := make([]map[string]interface{}, 0, len(files))

	for _, file := range files {
		result := map[string]interface{}{
			"id":        file.ID,
			"name":      file.Name,
			"createdAt": file.CreatedAt,
			"isFolder":  file.IsFolder,
			"parent":    file.Parent,
		}

		if file.IsFolder {
			result["isFolder"] = true
			result["parent"] = file.Parent
			children := file.Children
			if children == nil {
				children = []primitive.ObjectID{}
			}
			result["children"] = children
		} else {
			result["size"] = file.Size
			result["mimeType"] = file.MimeType
			result["path"] = file.Path
			result["parent"] = file.Parent
		}

		response = append(response, result)
	}

	if len(response) == 0 {
		response = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func CreateFolder(c *gin.Context) {
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

	folderValue, ok := c.Get(middleware.FolderContextKey)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Folder info not found in context"})
		return
	}

	folder, ok := folderValue.(models.File)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder data type"})
		return
	}

	folder.IsFolder = true
	folder.OwnerID = objID
	folder.CreatedAt = time.Now()
	folder.ID = primitive.NewObjectID()

	_, err = config.FilesCollection.InsertOne(ctx, folder)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Database request timed out"})
		} else if mongo.IsDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "Folder with this name already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder"})
		}
		return
	}

	response := map[string]interface{}{
		"id":        folder.ID,
		"createdAt": folder.CreatedAt,
		"name":      folder.Name,
	}

	c.JSON(http.StatusCreated, response)
}

func MoveFile(c *gin.Context) {
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

	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	targetID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID format"})
		return
	}

	fileValue, ok := c.Get(middleware.UpdateFileContextKey)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File info not found in context"})
		return
	}

	file, ok := fileValue.(models.File)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder data type"})
		return
	}

	var existingFile models.File
	err = config.FilesCollection.FindOne(ctx, bson.M{"_id": targetID, "ownerId": objUserID}).Decode(&existingFile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		}
		return
	}

	if file.Parent != nil && file.Parent.Hex() == targetID.Hex() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot move to itself"})
		return
	}

	updateTarget := bson.M{"$set": bson.M{"parent": file.Parent}}
	if result, err := config.FilesCollection.UpdateOne(ctx, bson.M{"_id": targetID}, updateTarget); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update target: " + err.Error()})
		return
	} else if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target file not found for update"})
		return
	}

	if existingFile.Parent != nil {
		updateOldParent := bson.M{"$pull": bson.M{"children": targetID}}
		if result, err := config.FilesCollection.UpdateOne(ctx, bson.M{"_id": existingFile.Parent}, updateOldParent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update old parent: " + err.Error()})
			return
		} else if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Old parent not found for update"})
			return
		}
	}

	if file.Parent != nil {
		parentID := *file.Parent
		fmt.Printf("Updating new parent with ID: %s, adding child: %s\n", parentID.Hex(), targetID.Hex())

		initChildren := bson.M{"$set": bson.M{"children": []primitive.ObjectID{}}}
		if _, err := config.FilesCollection.UpdateOne(ctx, bson.M{"_id": parentID, "children": nil}, initChildren); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize children: " + err.Error()})
			return
		}

		updateNewParent := bson.M{"$addToSet": bson.M{"children": targetID}}
		if result, err := config.FilesCollection.UpdateOne(ctx, bson.M{"_id": parentID}, updateNewParent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update new parent: " + err.Error()})
			return
		} else if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "New parent not found for update"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "File moved successfully"})
}

func DeleteFolder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
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

	folderID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID format"})
		return
	}

	var folder models.File
	if err := config.FilesCollection.FindOne(ctx, bson.M{"_id": folderID, "ownerId": objID}).Decode(&folder); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get folder: " + err.Error()})
		}
		return
	}

	if !folder.IsFolder {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Specified ID is not a folder"})
		return
	}

	totalSizeReduction, err := deleteFolderRecursive(ctx, folderID, userID, folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete folder contents: " + err.Error()})
		return
	}

	if result, err := config.FilesCollection.DeleteOne(ctx, bson.M{"_id": folderID, "ownerId": objID}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete folder: " + err.Error()})
		return
	} else if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found for deletion"})
		return
	}

	if folder.Parent != nil {
		updateParent := bson.M{"$pull": bson.M{"children": folderID}}
		if result, err := config.FilesCollection.UpdateOne(ctx, bson.M{"_id": *folder.Parent}, updateParent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update parent: " + err.Error()})
			return
		} else if result.MatchedCount == 0 {
			log.Printf("Parent %s not found for update, possibly already deleted", folder.Parent.Hex())
		}
	}

	if totalSizeReduction > 0 {
		if result, err := config.UserCollection.UpdateOne(
			ctx,
			bson.M{"_id": objID},
			bson.M{"$inc": bson.M{"usedStorage": -totalSizeReduction}},
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update storage: " + err.Error()})
			return
		} else if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found for storage update"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Folder and its contents deleted successfully",
		"sizeReduced": totalSizeReduction,
	})
}

func deleteFolderRecursive(ctx context.Context, folderID primitive.ObjectID, userID string, folder models.File) (int64, error) {
	var totalSizeReduction int64 = 0

	for _, childID := range folder.Children {
		var child models.File
		if err := config.FilesCollection.FindOne(ctx, bson.M{"_id": childID}).Decode(&child); err != nil {
			if err == mongo.ErrNoDocuments {
				log.Printf("Child %s not found, skipping", childID.Hex())
				continue
			}
			return 0, fmt.Errorf("failed to get child %s: %v", childID.Hex(), err)
		}

		if child.IsFolder {
			size, err := deleteFolderRecursive(ctx, childID, userID, child)
			if err != nil {
				return 0, fmt.Errorf("failed to delete subfolder %s: %v", childID.Hex(), err)
			}
			totalSizeReduction += size

			if result, err := config.FilesCollection.DeleteOne(ctx, bson.M{"_id": childID}); err != nil {
				return 0, fmt.Errorf("failed to delete subfolder %s: %v", childID.Hex(), err)
			} else if result.DeletedCount == 0 {
				log.Printf("Subfolder %s not found for deletion, possibly already deleted", childID.Hex())
			}
		} else {
			fileSizeMB := int64((child.Size + (1<<20 - 1)) / (1 << 20))
			if fileSizeMB < 1 {
				fileSizeMB = 1
			}
			totalSizeReduction += fileSizeMB

			if result, err := config.FilesCollection.DeleteOne(ctx, bson.M{"_id": childID}); err != nil {
				return 0, fmt.Errorf("failed to delete file %s: %v", childID.Hex(), err)
			} else if result.DeletedCount == 0 {
				log.Printf("File %s not found for deletion, possibly already deleted", childID.Hex())
			}

			filePath := fmt.Sprintf("uploads/%s/%s", userID, childID.Hex())
			if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
				log.Printf("Failed to delete file from disk %s: %v", filePath, err)
			}
		}
	}

	return totalSizeReduction, nil
}
