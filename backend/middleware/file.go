package middleware

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const maxFileSize = 5 << 20 // fix the size later

var forbiddenMimeTypes = map[string]bool{
	"application/x-msdownload":      true,
	"application/x-sh":              true,
	"application/x-php":             true,
	"application/javascript":        true,
	"application/x-bat":             true,
	"application/x-msdos-program":   true,
	"application/vnd.ms-powerpoint": true,
	"application/xhtml+xml":         true,
	"application/octet-stream":      true,
}

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

func MineFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "File is required"})
		return
	}

	defer file.Close()

	if isHasForbidden(header.Filename) {
		c.AbortWithStatusJSON(400, gin.H{"error": "File extension not allowed"})
		return
	}

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid file"})
		return
	}

	file.Seek(0, 0)

	mineType := http.DetectContentType(buffer)

	if forbiddenMimeTypes[mineType] {
		c.AbortWithStatusJSON(400, gin.H{"error": "File type not allowed"})
		return
	}

	c.Next()
}

func isHasForbidden(fileName string) bool {
	forbidden := []string{".exe", ".js", ".php", ".sh", ".bat", ".com", ".html", ".svg"}
	ext := strings.ToLower(filepath.Ext(fileName))

	for _, file := range forbidden {
		if ext == file {
			return true
		}
	}
	return false
}
