# EDNIC Backend

Golang API service for EDNIC.

## Local Setup

1. Prepare MySQL locally.

Use your existing local MySQL server, then create the local database if it does not exist:

```sql
CREATE DATABASE ednic_local CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. Start MinIO with Docker:

```text
docker compose up -d
```

3. Copy environment file:

```text
cp .env.example .env
```

4. Update `DATABASE_DSN` in `.env`.

5. Run API:

```text
go run ./cmd/api
```

## Database Migration

Run migrations from the repository root:

```text
go run ./cmd/migrate
```

## Bootstrap Super Admin

Set these values in `.env`:

```text
SUPER_ADMIN_NAME=
SUPER_ADMIN_USERNAME=
SUPER_ADMIN_EMAIL=
SUPER_ADMIN_PASSWORD=
```

Then run:

```text
go run ./cmd/bootstrap-admin
```

## Health Endpoints

```text
GET /health
GET /api/v1/health/db
```

`/health` verifies the API process is running.

`/api/v1/health/db` verifies the API can connect to MySQL.

## Auth Endpoints

```text
POST /api/v1/auth/register
POST /api/v1/auth/login
GET  /api/v1/auth/me
GET  /api/v1/admin/health
```

Protected endpoints use:

```text
Authorization: Bearer <token>
```

## Account Endpoints

```text
GET   /api/v1/account/profile
PATCH /api/v1/account/profile
PATCH /api/v1/account/password
```

## Product Endpoints

Public:

```text
GET /api/v1/products
GET /api/v1/products/:slug
GET /api/v1/product-categories
```

Admin:

```text
GET    /api/v1/admin/products
POST   /api/v1/admin/products
GET    /api/v1/admin/products/:id
PUT    /api/v1/admin/products/:id
DELETE /api/v1/admin/products/:id

GET    /api/v1/admin/product-categories
POST   /api/v1/admin/product-categories
PUT    /api/v1/admin/product-categories/:id
DELETE /api/v1/admin/product-categories/:id
```

## Swagger

Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

Regenerate Swagger docs after changing API annotations:

```text
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

## Local MinIO

MinIO API:

```text
http://localhost:9000
```

MinIO console:

```text
http://localhost:9001
```

Use `MINIO_ROOT_USER` and `MINIO_ROOT_PASSWORD` from `.env`.
