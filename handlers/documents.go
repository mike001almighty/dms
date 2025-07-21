package handlers

import (
	"log"

	"net/http"

	"dms/database"
	"dms/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateDocument(c *gin.Context) {
	log.Println("Creating document")
	var doc models.Document

	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := doc.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document"})
		return
	}
	log.Println("Document created successfully with id: ", doc.ID, "and title: ", doc.Title)
	c.JSON(http.StatusCreated, doc)
}

func GetDocument(c *gin.Context) {
	log.Println("Getting document with id: ", c.Param("id"))
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}
	log.Println("Document with id: ", id, "found")
	doc, err := models.GetDocumentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

func UpdateDocument(c *gin.Context) {
	log.Println("Updating document with id: ", c.Param("id"))
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
	log.Println("Deleting document with id: ", c.Param("id"))
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	err = models.DeleteDocumentByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

func BasicHealthCheck(c *gin.Context) {
	log.Println("Checking basic health")
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "dms",
	})
}

func DetailedHealthCheck(c *gin.Context) {
	log.Println("Checking detailed health")
	health := gin.H{
		"status":  "healthy",
		"service": "dms",
		"checks":  gin.H{},
	}

	// Check database connectivity
	if err := database.DB.Ping(c.Request.Context()); err != nil {
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

	health["checks"] = gin.H{
		"database": gin.H{
			"status": "healthy",
		},
	}
	c.JSON(http.StatusOK, health)
}
