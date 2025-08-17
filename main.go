package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"media-server/config"
	"media-server/handlers"
)

func main() {
	cfg := config.Load()

	router := gin.Default()


	router.POST("/upload", handlers.UploadFile(cfg))
	router.GET("/files", handlers.ListFiles(cfg))
	router.DELETE("/files/:name", handlers.DeleteFile(cfg))
	
	
	router.Static("/public", cfg.UploadDir)

	log.Printf("Starting server on port %s, upload dir: %s", cfg.Port, cfg.UploadDir)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
