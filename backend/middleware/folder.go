package middleware

import (
	"backend/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const FolderContextKey = "folderKey"

func FolderMiddlware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var folder models.File

		if err := c.ShouldBindJSON(&folder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Invalid JSON": err.Error()})
			c.Abort()
			return
		}

		if len(folder.Name) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Folder name must be at least 2 characters"})
			c.Abort()
			return
		}

		if strings.TrimSpace(folder.Name) != folder.Name || strings.HasPrefix(folder.Name, ".") || strings.HasSuffix(folder.Name, ".") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Folder name cannot start or end with space or dot"})
			c.Abort()
			return
		}

		c.Set(FolderContextKey, folder)

		c.Next()
	}
}
