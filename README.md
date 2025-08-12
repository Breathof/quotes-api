# Quotes API

A RESTful API for managing quotes and authors, built with Go, PostgreSQL, and Docker.

## Features

- **RESTful API** with JSON responses
- **PostgreSQL** database with migrations
- **Docker** and **Docker Compose** for easy deployment
- **Structured logging** with zerolog
- **Configuration via environment variables**
- **Health checks** (liveness and readiness probes)
- **Graceful shutdown**
- **Generated SQL queries** with sqlc for type safety
- **CI/CD** with GitHub Actions
- **API versioning** (v1)
- **Pagination** support
- **Search** functionality
- **CORS** support

## Architecture

The project follows a clean architecture pattern with the following layers:

```
cmd/server/          # Application entrypoint
internal/
  ├── api/          # HTTP handlers and routing
  ├── config/       # Configuration management
  ├── logger/       # Logging setup
  ├── repository/   # Data access layer
  └── service/      # Business logic layer
migrations/         # Database migrations
```

## API Endpoints

### Health Checks
- `GET /healthz` - Liveness probe
- `GET /readyz` - Readiness probe

### Authors
- `GET /api/v1/authors` - List all authors (paginated)
- `POST /api/v1/authors` - Create a new author
- `GET /api/v1/authors/{id}` - Get author by ID
- `PUT /api/v1/authors/{id}` - Update author
- `DELETE /api/v1/authors/{id}` - Delete author
- `GET /api/v1/authors/search?q={query}` - Search authors by name

### Quotes
- `GET /api/v1/quotes` - List all quotes (paginated)
- `POST /api/v1/quotes` - Create a new quote
- `GET /api/v1/quotes/{id}` - Get quote by ID
- `PUT /api/v1/quotes/{id}` - Update quote
- `DELETE /api/v1/quotes/{id}` - Delete quote
- `GET /api/v1/quotes/search?q={query}` - Search quotes by content
- `GET /api/v1/quotes/random` - Get a random quote
- `GET /api/v1/quotes?author_id={id}` - List quotes by author

## Getting Started

### Prerequisites

- Go 1.22 or higher
- Docker and Docker Compose
- PostgreSQL (if running locally)

### Running with Docker Compose

1. Clone the repository:
```bash
git clone https://github.com/igferreira/quotes-api.git
cd quotes-api
```

2. Start the services:
```bash
docker-compose up -d
```

The API will be available at `http://localhost:8080`

### Running Locally

1. Install dependencies:
```bash
go mod download
```

2. Set up the database:
```bash
# Start PostgreSQL (if not using Docker)
# Create a database named 'quotes_db'
```

3. Set environment variables:
```bash
export DB_PASSWORD=your_password
export DB_HOST=localhost
export DB_USER=postgres
export DB_NAME=quotes_db
```

4. Run migrations:
```bash
make migrate-up
```

5. Generate SQL code:
```bash
make sqlc-generate
```

6. Run the application:
```bash
make run
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | *required* |
| `DB_NAME` | Database name | `quotes` |
| `DB_SSL_MODE` | Database SSL mode | `disable` |
| `LOG_LEVEL` | Log level (debug, info, warn, error) | `info` |
| `LOG_JSON` | Output logs in JSON format | `false` |
| `ENVIRONMENT` | Environment (development, production) | `development` |

## Development

### Install development tools:
```bash
make dev-setup
```

### Run tests:
```bash
make test
```

### Run linters:
```bash
make lint
```

### Format code:
```bash
make fmt
```

## API Examples

### Create an author:
```bash
curl -X POST http://localhost:8080/api/v1/authors \
  -H "Content-Type: application/json" \
  -d '{"name": "Albert Einstein", "bio": "Theoretical physicist"}'
```

### Create a quote:
```bash
curl -X POST http://localhost:8080/api/v1/quotes \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Imagination is more important than knowledge.",
    "author_id": 1,
    "source": "Interview (1929)",
    "tags": ["imagination", "knowledge"]
  }'
```

### Get a random quote:
```bash
curl http://localhost:8080/api/v1/quotes/random
```

### Search quotes:
```bash
curl "http://localhost:8080/api/v1/quotes/search?q=imagination&limit=10"
```

## License

This project is licensed under the MIT License.