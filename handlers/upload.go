package handlers

import (
	"io"
	"log"
	"media-server/config"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// handlers/upload.go (обновлённый)
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

		ext := strings.ToLower(filepath.Ext(file.Filename))
		fileID := uuid.New().String() + ext
		path := filepath.Join(cfg.UploadDir, fileID)

		if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
			log.Println("mkdir error:", err)
		}

		// Перемещаем указатель файла в начало
		src.Seek(0, 0)
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
			return
		}

		// Собираем дополнительные поля из формы (кроме "file")
		formValues := make(gin.H)
		for key, values := range c.Request.PostForm {
			if key != "file" && len(values) > 0 {
				formValues[key] = values[0]
			}
		}

		// Добавляем информацию о файле
		result := gin.H{
			"id":  fileID,
			"url": cfg.Domain + cfg.PublicPath + "/" + fileID,
		}
		for k, v := range formValues {
			result[k] = v
		}

		c.JSON(http.StatusOK, result)
	}
}
