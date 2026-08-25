# Temporary File Sharing Service

A small backend service built with Go that lets users upload files and share them through a temporary download link.

The idea is simple: upload a file, get an ID, share that ID, and the file automatically expires after a certain amount of time.

## Features
- Upload files through an API
- Download uploaded files using a unique file ID
- Files automatically expire after 24 hours
- Expired files are removed automatically by a background worker
- Maximum upload size is configurable
- Basic file type validation
- PostgreSQL for storing file metadata
- Local storage for the actual files
- Docker for PostgreSQL
- Graceful shutdown for the background worker
  
## Tech Stack

- Go
- Gin
- PostgreSQL
- pgx
- Docker
- Docker Compose

## Running the Project
1. Clone the repository
` git clone github.com/Golukpal/Temporary-File-Sharing- `
cd temporary-file-sharing
2. Start PostgreSQL

The project uses Docker for PostgreSQL.

` docker compose up -d `

Check that the container is running:

`docker ps`
3. Create the database table

Run the SQL from:

`migrations/001_create_files.sql`

against the file_sharing database.

For example:

`docker exec -i temporary_file_postgres \
  psql -U postgres -d file_sharing \
  < migrations/001_create_files.sql`
  
4. Install Go dependencies
`go mod tidy`
5. Start the API
`go run ./cmd/server`

The server will start on:

http://localhost:8080

You can check it with:

curl http://localhost:8080/health

## Architecture

The project follows a simple layered structure:

``` HTTP Request
     ↓
Handler
     ↓
Service
     ↓
Repository
     ↓
PostgreSQL 

File storage is handled separately: 

Service
   ↓
storage/
   ↓
Actual file

There is also a background cleanup worker: 

Cleanup Worker
      ↓
Find expired files
      ↓
Delete physical file
      ↓
Delete database record
