package models

import (
	"context"
	"time"

	"dms/database"
)

type Document struct {
	ID          int       `json:"id"`
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
