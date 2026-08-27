# Document Management System (DMS)

A simple multi-tenant document API built with Go, PostgreSQL, and Keycloak.

## Features

- Login via Keycloak (password grant) returning a JWT
- Create, list, get, and delete documents
- Store file content as base64-encoded data in PostgreSQL
- Tenant scoping derived from the JWT (documents are filtered by `tenant_id`)

## Prerequisites

- Docker and Docker Compose
- Go 1.23+ (for local development outside Docker)

## Running the Application

```bash
docker compose up -d --build
```

Services:

| Service   | URL / port                          |
|-----------|-------------------------------------|
| API       | `http://localhost:8080`             |
| Keycloak  | `http://localhost:8082` (admin UI: `/admin`, user `admin` / `admin`) |
| Postgres  | `localhost:5433`                    |

Keycloak realm `dms` (client `dms-service`, user `dms_admin`) is imported automatically from `dms-realm.json` on first start (empty Keycloak volume).

### Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"dms_admin","password":"dms_admin_password"}'
```

Use the returned `access_token` as `Authorization: Bearer <token>` on document endpoints.

## API Endpoints

| Method   | Path              | Auth   | Description              |
|----------|-------------------|--------|--------------------------|
| `POST`   | `/login`          | No     | Exchange credentials for JWT |
| `POST`   | `/documents`      | Bearer | Create a document        |
| `GET`    | `/documents`      | Bearer | List documents for tenant |
| `GET`    | `/documents/:id`  | Bearer | Get one document         |
| `DELETE` | `/documents/:id`  | Bearer | Delete a document        |
| `GET`    | `/health`         | No     | Liveness                 |
| `GET`    | `/health/detailed`| No     | Liveness + DB check      |

### Create document

`POST /documents` — `Content-Type: application/json`

```json
{
  "title": "My Document",
  "extension": "pdf",
  "description": "A sample PDF document",
  "content": "base64-encoded-file-content"
}
```

Response includes `id` (UUID), `tenant_id`, timestamps, and the fields above.

### List documents

`GET /documents` — returns an array of documents for the caller's tenant.

## Local development (app only)

```bash
docker compose up -d db keycloak

export DB_HOST=localhost
export DB_PORT=5433
export DB_USER=dms_user
export DB_PASSWORD=dms_password
export DB_NAME=dms
export KEYCLOAK_URL=http://localhost:8082
export KEYCLOAK_REALM=dms
export KEYCLOAK_CLIENT_ID=dms-service
export KEYCLOAK_CLIENT_SECRET=your-service-secret-key

go run main.go
```

API listens on `:8080`.

## Project structure

```
dms/
├── main.go
├── go.mod / go.sum
├── docker-compose.yml
├── Dockerfile
├── init.sql
├── dms-realm.json          # Keycloak realm import
├── keycloak_setup_guide.md # Manual Keycloak steps (optional; import covers this)
├── auth/                   # JWT middleware + Keycloak client
├── handlers/
├── models/
└── database/
```
