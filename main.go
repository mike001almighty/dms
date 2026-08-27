package main

import (
	"log"

	"dms/auth"
	"dms/database"
	"dms/handlers"

	"github.com/gin-gonic/gin"
)


func main() {
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.DB.Close()

	r := gin.Default()

	// Login endpoint (no authentication required)
	r.POST("/login", handlers.Login)

	authMw := auth.JWTAuthMiddleware()
	canReadWrite := auth.RequireAnyRole("DMS_USER", "DMS_ADMIN")
	canDelete := auth.RequireRole("DMS_ADMIN")

	r.POST("/documents", authMw, canReadWrite, handlers.CreateDocument)
	r.GET("/documents", authMw, canReadWrite, handlers.ListDocuments)
	r.GET("/documents/:id", authMw, canReadWrite, handlers.GetDocument)
	r.DELETE("/documents/:id", authMw, canDelete, handlers.DeleteDocument)

	// Health endpoints don't need tenant filtering
	r.GET("/health", handlers.BasicHealthCheck)
	r.GET("/health/detailed", handlers.DetailedHealthCheck)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
