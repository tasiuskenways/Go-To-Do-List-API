# 🚀 Go To-Do List API

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/tasiuskenways/todolistapi.svg)](https://pkg.go.dev/github.com/tasiuskenways/todolistapi)

## ℹ️ About

This project is a personal portfolio implementation of a Todo List API, inspired by the [Todo List API project](https://roadmap.sh/projects/todo-list-api) from roadmap.sh. It's built to demonstrate backend development skills and best practices in Go.

A production-ready, secure, and scalable To-Do List API built with Go, Fiber, GORM, and Redis. This API provides JWT-based authentication, hybrid encryption, and full CRUD operations for managing todo items with enterprise-grade security and performance.

## 🚀 Quick Start

Get up and running in under 5 minutes:

```bash
# 1. Clone and enter the repository
git clone https://github.com/tasiuskenways/Go-To-Do-List-API.git
cd Go-To-Do-List-API

# 2. Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# 3. Install dependencies
go mod download

# 4. Generate necessary keys
go run cmd/generate_jwt_key/main.go
go run cmd/generate_keys/main.go -output ./keys

# 5. Start the database and cache
docker-compose up -d postgres redis

# 6. Run migrations
go run cmd/migrate/main.go

# 7. Start the server
go run cmd/main.go
```

Your API is now running at `http://localhost:3000` 🎉

## Features

- 🔐 JWT-based authentication with refresh tokens
- 🔒 Hybrid encryption for sensitive data
- 📝 CRUD operations for todo items
- 🚀 High performance with Fiber framework
- 🗄️ PostgreSQL for data persistence
- 🎯 Redis for caching and token management
- 🐳 Docker support
- 🔄 Database migrations
- ✅ Input validation
- 📊 Structured logging

## Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Redis
- Make (optional)

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/todolistapi.git
cd todolistapi
```

### 2. Set up environment variables

Copy the example environment file and update the values:

```bash
cp .env.example .env
```

### 3. Install dependencies

```bash
go mod download
```

### 4. Generate JWT keys

```bash
go run cmd/generate_jwt_key/main.go
```

### 5. Generate encryption keys

```bash
go run cmd/generate_keys/main.go -output ./keys
```

### 6. Run database migrations

```bash
go run cmd/migrate/main.go
```

### 7. Start the server

```bash
go run cmd/main.go
```

The API will be available at `http://localhost:3000`

## 📚 API Documentation

### Response Format
All responses follow this format:
```json
{
   "requestId": "unique_request_id",
   "success": true,
   "message": "Operation successful",
   "data": {},
   "errors": null
}
```

### Error Handling
Common error responses:
- `400 Bad Request` - Invalid input data
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Project Structure

```
.
├── cmd/                  # Main applications
│   ├── main.go           # Entry point
│   └── migrate/          # Database migration
└── internal/             # Private application code
    ├── application/      # Application services
    ├── config/           # Configuration
    ├── domain/           # Domain models and interfaces
    └── interfaces/       # Interface adapters (HTTP, gRPC, etc.)
        └── http/         # HTTP handlers and routes
```

## Available Commands

The `/cmd` directory contains several utility commands:

### 1. Main Application
```bash
# Run the main application
go run cmd/main.go
```

### 2. Database Migrations
```bash
# Run database migrations
go run cmd/migrate/main.go
```

### 3. Key Generation
```bash
# Generate JWT secret key
go run cmd/generate_jwt_key/main.go

# Generate RSA key pair for hybrid encryption
go run cmd/generate_keys/main.go -output ./keys
```

### 4. Hybrid Encryption (Development)
```bash
# Encrypt data
go run cmd/hybrid_encryption/main.go -action encrypt -input '{"key":"value"}'

# Decrypt data
go run cmd/hybrid_encryption/main.go -action decrypt -input "<encrypted_data>"
```

## Development

### Running tests

```bash
go test ./... -v
```

### Building the application

```bash
go build -o todolistapi cmd/main.go
```

### Running with Docker

```bash
docker-compose up --build
```

## 🔒 Security

### Authentication & Authorization
- JWT-based authentication with access and refresh tokens
- Short-lived access tokens (15 minutes by default)
- Long-lived refresh tokens (30 days by default)
- Token blacklisting on logout
- Role-based access control (RBAC) ready

### Data Protection
- Hybrid encryption (RSA + AES) for sensitive data
- Password hashing using bcrypt with work factor 12
- Request/response encryption for sensitive endpoints
- SQL injection prevention through GORM
- XSS protection headers
- CORS configuration with secure defaults

### Best Practices
- Environment-based configuration
- Secure defaults in code
- No sensitive data in logs
- Rate limiting (recommended to implement at the proxy level)

## 🤝 Contributing

We welcome contributions! Here's how to get started:

1. Fork the repository and create your feature branch
   ```bash
   git checkout -b feature/amazing-feature
   ```
2. Make your changes and ensure tests pass
   ```bash
   go test ./...
   ```
3. Commit your changes with a descriptive message
   ```bash
   git commit -m 'feat: add amazing feature'
   ```
4. Push to your fork
   ```bash
   git push origin feature/amazing-feature
   ```
5. Open a Pull Request with a clear description

### Development Workflow
- Write tests for new features
- Update documentation when adding new endpoints
- Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification
- Keep the code style consistent (we use `gofmt` and `golint`)

## 🚀 Deployment

### Docker Deployment
```bash
docker-compose up --build -d
```

### Environment Variables
Configure the application using these environment variables:

```env
# Application
APP_ENV=development
APP_PORT=3000

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todolist

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_PRIVATE_KEY=your_jwt_private_key
JWT_EXPIRATION=15m
JWT_REFRESH_EXPIRATION=720h
```

## 🐛 Troubleshooting

### Common Issues
1. **Database Connection Issues**
   - Ensure PostgreSQL is running and accessible
   - Check database credentials in `.env`
   - Verify the database exists and migrations have run

2. **JWT Errors**
   - Verify `JWT_PRIVATE_KEY` is set and consistent
   - Check token expiration times
   - Ensure system clock is synchronized (for JWT validation)

3. **Encryption Issues**
   - Ensure keys are generated and accessible
   - Check file permissions for key files
   - Verify the correct key paths in configuration

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fiber](https://github.com/gofiber/fiber) - Fast HTTP framework for Go
- [GORM](https://gorm.io/) - The fantastic ORM library for Golang
- [Redis](https://redis.io/) - In-memory data structure store
- [JWT](https://jwt.io/) - JSON Web Tokens
- [Docker](https://www.docker.com/) - Containerization platform
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - Password hashing