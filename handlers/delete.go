package handlers

import (
	"media-server/config"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func DeleteFile(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing filename"})
			return
		}

		// защита от directory traversal ("../")
		cleanName := filepath.Base(name)
		path := filepath.Join(cfg.UploadDir, cleanName)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}

		if err := os.Remove(path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"deleted": true, "file": cleanName})
	}
}
