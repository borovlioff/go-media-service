// handlers/uploadBlob.go
package handlers

import (
	"encoding/base64"
	"log"
	"media-server/config"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadBlobRequest struct {
	File string `json:"file"` // base64-encoded blob
}

func UploadBlob(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UploadBlobRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}

		if req.File == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing file in request"})
			return
		}

		// Поддержка data URL: `data:image/png;base64,...`
		prefix := "data:"
		if strings.HasPrefix(req.File, prefix) {
			commaIdx := strings.Index(req.File, ",")
			if commaIdx > 0 {
				req.File = req.File[commaIdx+1:]
			}
		}

		data, err := base64.StdEncoding.DecodeString(req.File)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid base64"})
			return
		}

		if int64(len(data)) > cfg.FileMaxSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file too large"})
			return
		}

		// Определение MIME-типа по первым байтам
		mimeType := http.DetectContentType(data)
		if !cfg.AllowedMIMEs[mimeType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed", "mime": mimeType})
			return
		}

		ext := ".bin"
		if idx := strings.Index(mimeType, "/"); idx != -1 {
			parts := strings.Split(mimeType, "/")
			if parts[1] != "" {
				ext = "." + parts[1]
				// Унификация расширений
				switch ext {
				case ".jpeg":
					ext = ".jpg"
				case ".mpeg":
					ext = ".mp4"
				}
			}
		}

		newName := uuid.New().String() + ext
		path := filepath.Join(cfg.UploadDir, newName)

		if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
			log.Println("mkdir error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create upload dir"})
			return
		}

		file, err := os.Create(path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create file"})
			return
		}
		defer file.Close()

		_, err = file.Write(data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":  newName,
			"url": cfg.Domain + cfg.PublicPath + "/" + newName,
		})
	}
}
