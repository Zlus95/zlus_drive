package middleware

import (
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

const StorageLimitContextKey = "storageKey"

func StorageLimitMiddlware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Invalid JSON": err.Error()})
			c.Abort()
			return
		}

		if user.StorageLimit < user.UsedStorage {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Storage Limit must be greater than Used Storage"})
			c.Abort()
			return
		}

		if user.StorageLimit > 15 * 1024 * 1024 * 1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Storage Limit must be less than 15GB"})
			c.Abort()
			return
		}

		c.Set(StorageLimitContextKey, user.StorageLimit)

		c.Next()
	}
}
