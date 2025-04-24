package middleware

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const maxFileSize = 50 << 20

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

var forbidden = map[string]bool{
	".exe":  true,
	".js":   true,
	".php":  true,
	".sh":   true,
	".bat":  true,
	".com":  true,
	".html": true,
}

type FileInfo struct {
	File     multipart.File
	Header   *multipart.FileHeader
	MimeType string
	Size     int64
	Filename string
}

const FileContextKey = "fileInfo"

func FileMiddlware(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	if header.Size > maxFileSize {
		c.AbortWithStatusJSON(400, gin.H{"error": fmt.Sprintf("File size exceeds the limit (%dMB)", maxFileSize>>20)})
		return
	}

	if isHasForbidden(header.Filename) {
		c.AbortWithStatusJSON(400, gin.H{"error": "File extension not allowed"})
		return
	}

	mimeType, err := detectMimeType(file)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid file"})
		return
	}

	if forbiddenMimeTypes[mimeType] {
		c.AbortWithStatusJSON(400, gin.H{"error": "File type not allowed"})
		return
	}

	fileInfo := FileInfo{
		File:     file,
		Header:   header,
		MimeType: mimeType,
		Size:     header.Size,
		Filename: header.Filename,
	}

	c.Set(FileContextKey, fileInfo)

	c.Next()

}

func isHasForbidden(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	return forbidden[ext]
}

func detectMimeType(file multipart.File) (string, error) {
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return "", err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}
	return http.DetectContentType(buffer), nil
}
