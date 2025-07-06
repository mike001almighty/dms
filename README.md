# Document Management System (DMS)

A simple document management system built with Go and PostgreSQL.

## Features

- Upload documents with metadata (title, extension, description)
- Store file content as base64 encoded data
- PostgreSQL database for persistent storage
- RESTful API endpoint

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)

### Running the Application

1. Start the services:
```bash
docker-compose up -d
```

2. The API will be available at `http://localhost:8080`

### API Endpoints

#### Upload Document
- **POST** `/documents`
- **Content-Type:** `application/json`

**Request Body:**
```json
{
    "title": "My Document",
    "extension": "pdf",
    "description": "A sample PDF document", 
    "content": "base64-encoded-file-content"
}
```

**Response:**
```json
{
    "id": 1,
    "title": "My Document",
    "extension": "pdf", 
    "description": "A sample PDF document",
    "content": "base64-encoded-file-content",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
}
```

## Development

### Local Development
```bash
# Install dependencies
go mod tidy

# Run PostgreSQL
docker-compose up -d db

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=dms_user
export DB_PASSWORD=dms_password
export DB_NAME=dms

# Run the application
go run main.go
```

### Project Structure
```
dms/
├── main.go              # Application entry point
├── go.mod               # Go module definition
├── docker-compose.yml   # Docker services configuration
├── Dockerfile           # Go application container
├── init.sql             # Database schema
├── handlers/
│   └── documents.go     # API handlers
├── models/
│   └── document.go      # Data models
└── database/
    └── connection.go    # Database connection
``` 