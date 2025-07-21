package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func Connect() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = pgx.Connect(context.Background(), connString)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	return nil
}
