package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"media-server/config"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"io"
)

func UploadFile(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
			return
		}

		if file.Size > cfg.FileMaxSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file too large"})
			return
		}

		// Проверка MIME-типа
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
			return
		}
		defer src.Close()

		buffer := make([]byte, 512)
		_, err = src.Read(buffer)
		if err != nil && err != io.EOF {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot read file"})
			return
		}

		mimeType := http.DetectContentType(buffer)
		if !cfg.AllowedMIMEs[mimeType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed", "mime": mimeType})
			return
		}

		// Сохраняем файл с уникальным именем
		ext := strings.ToLower(filepath.Ext(file.Filename))
		newName := uuid.New().String() + ext
		path := filepath.Join(cfg.UploadDir, newName)

		if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
			log.Println("mkdir error:", err)
		}

		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"file": "/public/" + newName})
	}
}
