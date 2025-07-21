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
	log.Println("Creating document, with tenant id: ", c.MustGet("tenant_id"))
	tenantID := c.MustGet("tenant_id").(string)
	var doc models.Document

	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set tenant ID from context
	doc.TenantID = tenantID

	if err := doc.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document"})
		return
	}
	log.Println("Document created successfully with id: ", doc.ID, "and title: ", doc.Title, "and tenant id: ", doc.TenantID)
	c.JSON(http.StatusCreated, doc)
}

func GetDocument(c *gin.Context) {
	log.Println("Getting document with id: ", c.Param("id"), "and tenant id: ", c.MustGet("tenant_id"))
	tenantID := c.MustGet("tenant_id").(string)
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}
	doc, err := models.GetDocumentByID(id, tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

func UpdateDocument(c *gin.Context) {
	log.Println("Updating document with id: ", c.Param("id"), "and tenant id: ", c.MustGet("tenant_id"))
	tenantID := c.MustGet("tenant_id").(string)
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

	// Check if document exists within tenant scope
	existingDoc, err := models.GetDocumentByID(id, tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Update fields with tenant context
	doc.ID = id
	doc.TenantID = tenantID
	doc.CreatedAt = existingDoc.CreatedAt

	if err := doc.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

func DeleteDocument(c *gin.Context) {
	log.Println("Deleting document with id: ", c.Param("id"), "and tenant id: ", c.MustGet("tenant_id"))
	tenantID := c.MustGet("tenant_id").(string)
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	err = models.DeleteDocumentByID(id, tenantID)
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
