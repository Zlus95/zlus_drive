package middleware

import "github.com/gin-gonic/gin"

const maxFileSize = 5 << 20 // fix the size later

func FileSizeMiddlWare(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	if header.Size > maxFileSize {
		c.AbortWithStatusJSON(400, gin.H{"error": "File size exceeds the limit (5MB)"})
		return
	}

	c.Next()
}
