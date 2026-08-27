package models

import (
	"context"
	"time"

	"dms/database"

	"github.com/google/uuid"
)

type Document struct {
	ID          uuid.UUID `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Title       string    `json:"title"`
	Extension   string    `json:"extension"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (d *Document) Save() error {
	query := `
		INSERT INTO documents (tenant_id, title, extension, description, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	row := database.DB.QueryRow(context.Background(), query, d.TenantID, d.Title, d.Extension, d.Description, d.Content)
	err := row.Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func GetDocumentByID(id uuid.UUID, tenantID string) (*Document, error) {
	query := `
		SELECT id, tenant_id, title, extension, description, content, created_at, updated_at
		FROM documents
		WHERE id = $1 AND tenant_id = $2`

	var doc Document
	row := database.DB.QueryRow(context.Background(), query, id, tenantID)
	err := row.Scan(&doc.ID, &doc.TenantID, &doc.Title, &doc.Extension, &doc.Description, &doc.Content, &doc.CreatedAt, &doc.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func ListDocumentsByTenant(tenantID string) ([]Document, error) {
	query := `
		SELECT id, tenant_id, title, extension, description, content, created_at, updated_at
		FROM documents
		WHERE tenant_id = $1
		ORDER BY created_at DESC`

	rows, err := database.DB.Query(context.Background(), query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	docs := make([]Document, 0)
	for rows.Next() {
		var doc Document
		if err := rows.Scan(&doc.ID, &doc.TenantID, &doc.Title, &doc.Extension, &doc.Description, &doc.Content, &doc.CreatedAt, &doc.UpdatedAt); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return docs, nil
}

func DeleteDocumentByID(id uuid.UUID, tenantID string) error {
	query := `DELETE FROM documents WHERE id = $1 AND tenant_id = $2`
	_, err := database.DB.Exec(context.Background(), query, id, tenantID)
	return err
}
