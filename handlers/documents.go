package handlers

import (
	"net/http"

	"dms/database"
	"dms/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func getRequestLogger(c *gin.Context) *logrus.Entry {
	requestID, _ := c.Get("request_id")
	return logrus.WithField("request_id", requestID)
}

func CreateDocument(c *gin.Context) {
	logger := getRequestLogger(c)
	var doc models.Document

	if err := c.ShouldBindJSON(&doc); err != nil {
		logger.WithError(err).Warn("Invalid JSON in create document request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.WithFields(logrus.Fields{
		"title":     doc.Title,
		"extension": doc.Extension,
	}).Info("Creating new document")

	if err := doc.Save(); err != nil {
		logger.WithError(err).Error("Failed to save document to database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document"})
		return
	}

	logger.WithField("document_id", doc.ID).Info("Document created successfully")
	c.JSON(http.StatusCreated, doc)
}

func GetDocument(c *gin.Context) {
	logger := getRequestLogger(c)
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"provided_id": idStr,
			"error":       err.Error(),
		}).Warn("Invalid UUID provided for get document")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	logger.WithField("document_id", id).Info("Retrieving document")

	doc, err := models.GetDocumentByID(id)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"document_id": id,
			"error":       err.Error(),
		}).Warn("Document not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	logger.WithField("document_id", id).Info("Document retrieved successfully")
	c.JSON(http.StatusOK, doc)
}

func UpdateDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var doc models.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if document exists
	existingDoc, err := models.GetDocumentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Update fields
	doc.ID = id
	doc.CreatedAt = existingDoc.CreatedAt

	if err := doc.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

func DeleteDocument(c *gin.Context) {
	logger := getRequestLogger(c)
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"provided_id": idStr,
			"error":       err.Error(),
		}).Warn("Invalid UUID provided for delete document")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	logger.WithField("document_id", id).Info("Deleting document")

	err = models.DeleteDocumentByID(id)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"document_id": id,
			"error":       err.Error(),
		}).Error("Failed to delete document from database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	logger.WithField("document_id", id).Info("Document deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

func BasicHealthCheck(c *gin.Context) {
	getRequestLogger(c).Info("Basic health check requested")
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "dms",
	})
}

func DetailedHealthCheck(c *gin.Context) {
	logger := getRequestLogger(c)
	logger.Info("Detailed health check requested")

	health := gin.H{
		"status":  "healthy",
		"service": "dms",
		"checks":  gin.H{},
	}

	// Check database connectivity
	if err := database.DB.Ping(c.Request.Context()); err != nil {
		logger.WithError(err).Error("Database health check failed")
		health["status"] = "unhealthy"
		health["checks"] = gin.H{
			"database": gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			},
		}
		c.JSON(http.StatusServiceUnavailable, health)
		return
	}

	logger.Info("All health checks passed")
	health["checks"] = gin.H{
		"database": gin.H{
			"status": "healthy",
		},
	}
	c.JSON(http.StatusOK, health)
}
