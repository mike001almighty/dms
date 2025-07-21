package main

import (
	"log"
	"net/http"

	"dms/database"
	"dms/handlers"

	"github.com/gin-gonic/gin"
)

func tenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "X-Tenant-ID header is required"})
			c.Abort()
			return
		}

		// Add tenant ID to context for use in handlers
		c.Set("tenant_id", tenantID)
		c.Next()
	}
}

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := gin.Default()

	// Apply tenant middleware to document endpoints
	r.POST("/documents", tenantMiddleware(), handlers.CreateDocument)
	r.GET("/documents/:id", tenantMiddleware(), handlers.GetDocument)
	r.DELETE("/documents/:id", tenantMiddleware(), handlers.DeleteDocument)

	// Health endpoints don't need tenant filtering
	r.GET("/health", handlers.BasicHealthCheck)
	r.GET("/health/detailed", handlers.DetailedHealthCheck)

	log.Println("Server starting on :8085")
	if err := r.Run(":8085"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
