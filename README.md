# TodoList API

A Go-based TodoList API with JWT authentication, Redis caching, and PostgreSQL database.

## CI/CD Setup

This project uses GitHub Actions for CI/CD. The workflow includes:

1. **Build and Test**: Builds the application and runs basic checks
2. **Docker Build & Push**: Builds and pushes a Docker image to Docker Hub
3. **Deploy**: Deploys the application to the staging/production server

### Required Secrets

Add these secrets to your GitHub repository settings:

- `DOCKERHUB_USERNAME`: Your Docker Hub username
- `DOCKERHUB_TOKEN`: Your Docker Hub access token
- `SSH_PRIVATE_KEY`: SSH private key for deployment
- `SSH_USER`: SSH username for the deployment server
- `SSH_HOST`: Deployment server hostname or IP

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```
# Database
DB_HOST=postgres_db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todolist

# JWT
JWT_SECRET=your_jwt_secret
JWT_EXPIRATION=24h

# Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=

# App
APP_PORT=3000
APP_ENV=development
```

## Development

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Make (optional)

### Running Locally

1. Start the development environment:
   ```bash
   docker-compose build --no-cache
   docker-compose up -d
   ```

## API Documentation

API documentation is available at `http://localhost:3000/swagger/index.html` when running locally.

## Code Quality

This project uses:

- CodeClimate for code quality monitoring
- GitHub Actions for CI/CD
- Dependabot for dependency updates

## License

MIT
