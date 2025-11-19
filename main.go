package main

import (
	"log"

	"media-server/config"
	"media-server/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	router := gin.Default()

	router.POST("/upload", handlers.UploadFile(cfg))
	router.POST("/upload-blob", handlers.UploadBlob(cfg))
	router.GET("/files", handlers.ListFiles(cfg))
	router.DELETE("/files/:name", handlers.DeleteFile(cfg))

	router.Static(cfg.PublicPath, cfg.UploadDir)

	log.Printf("Starting server on port %s, upload dir: %s, public path: %s",
		cfg.Port, cfg.UploadDir, cfg.PublicPath)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
