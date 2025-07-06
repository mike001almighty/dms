package models

import (
	"context"
	"time"

	"dms/database"

	"github.com/google/uuid"
)

type Document struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Extension   string    `json:"extension"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (d *Document) Save() error {
	query := `
		INSERT INTO documents (title, extension, description, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	row := database.DB.QueryRow(context.Background(), query, d.Title, d.Extension, d.Description, d.Content)
	err := row.Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func GetDocumentByID(id uuid.UUID) (*Document, error) {
	query := `
		SELECT id, title, extension, description, content, created_at, updated_at
		FROM documents
		WHERE id = $1`

	var doc Document
	row := database.DB.QueryRow(context.Background(), query, id)
	err := row.Scan(&doc.ID, &doc.Title, &doc.Extension, &doc.Description, &doc.Content, &doc.CreatedAt, &doc.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func DeleteDocumentByID(id uuid.UUID) error {
	query := `DELETE FROM documents WHERE id = $1`
	_, err := database.DB.Exec(context.Background(), query, id)
	return err
}
