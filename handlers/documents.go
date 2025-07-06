package handlers

import (
	"net/http"

	"dms/models"

	"github.com/gin-gonic/gin"
)

func CreateDocument(c *gin.Context) {
	var doc models.Document

	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := doc.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document"})
		return
	}

	c.JSON(http.StatusCreated, doc)
}
