package main

import (
	"log"

	"dms/auth"
	"dms/database"
	"dms/handlers"

	"github.com/gin-gonic/gin"
)

// JWT middleware is now handled by auth package

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := gin.Default()

	// Apply JWT authentication middleware to document endpoints
	r.POST("/documents", auth.JWTAuthMiddleware(), handlers.CreateDocument)
	r.GET("/documents/:id", auth.JWTAuthMiddleware(), handlers.GetDocument)
	r.DELETE("/documents/:id", auth.JWTAuthMiddleware(), handlers.DeleteDocument)

	// Health endpoints don't need tenant filtering
	r.GET("/health", handlers.BasicHealthCheck)
	r.GET("/health/detailed", handlers.DetailedHealthCheck)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
