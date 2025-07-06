package main

import (
	"log"

	"dms/database"
	"dms/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := gin.Default()
	r.POST("/documents", handlers.CreateDocument)
	r.GET("/documents/:id", handlers.GetDocument)
	r.DELETE("/documents/:id", handlers.DeleteDocument)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
