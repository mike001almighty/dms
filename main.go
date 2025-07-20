package main

import (
	"time"

	"dms/database"
	"dms/handlers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func init() {
	// Configure structured logging
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logrus.SetLevel(logrus.InfoLevel)
}

func requestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add request ID to context and response headers
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		start := time.Now()

		// Log request start
		logrus.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"user_ip":    c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Info("Request started")

		c.Next()

		// Log request completion
		duration := time.Since(start)
		logrus.WithFields(logrus.Fields{
			"request_id":  requestID,
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status_code": c.Writer.Status(),
			"duration_ms": duration.Milliseconds(),
			"user_ip":     c.ClientIP(),
		}).Info("Request completed")
	}
}

func main() {
	if err := database.Connect(); err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}

	r := gin.Default()

	// Add request logging middleware
	r.Use(requestLoggingMiddleware())

	r.POST("/documents", handlers.CreateDocument)
	r.GET("/documents/:id", handlers.GetDocument)
	r.DELETE("/documents/:id", handlers.DeleteDocument)

	r.GET("/health", handlers.BasicHealthCheck)
	r.GET("/health/detailed", handlers.DetailedHealthCheck)

	logrus.WithField("port", "8085").Info("Server starting")
	if err := r.Run(":8085"); err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}
}
