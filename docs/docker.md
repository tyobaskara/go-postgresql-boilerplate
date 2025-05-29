# Docker Documentation

This document provides comprehensive information about the Docker setup for the Jeki Backend project.

## Overview

The project uses Docker for containerization and includes multiple services:
- Main application (Go)
- PostgreSQL database
- Swagger UI for API documentation

## Docker Files

### Dockerfile
The project includes two Dockerfile variants:

1. `Dockerfile` - Production build
   - Uses multi-stage build to minimize final image size
   - Based on `golang:1.22-alpine` for building
   - Final image based on `alpine:3.19`
   - Exposes port 8080

2. `Dockerfile.dev` - Development build
   - Optimized for development workflow
   - Includes hot-reloading capabilities

## Docker Compose

The `docker-compose.yml` file defines three services:

### 1. App Service
- Builds from either `Dockerfile` or `Dockerfile.dev` (configurable via `DOCKERFILE` env var)
- Exposes port 8080
- Mounts the current directory for development
- Uses environment variables from `.env.{ENV}` file
- Depends on PostgreSQL service

### 2. PostgreSQL Service
- Uses `postgres:16-alpine` image
- Exposes port 5432
- Environment variables:
  - `POSTGRES_USER` (default: postgres)
  - `POSTGRES_PASSWORD` (default: postgres)
  - `POSTGRES_DB` (default: jeki)
- Includes health checks
- Initializes database using `scripts/init.sql`

### 3. Swagger UI Service
- Uses `swaggerapi/swagger-ui` image
- Exposes port 8081
- Mounts Swagger documentation from `docs/swagger`

## Port Mapping and Expose

### Understanding Port Mapping

In Docker, there are two important concepts related to ports:
1. `EXPOSE` in Dockerfile
2. Port mapping in docker-compose.yml

#### Port Mapping Format
```yaml
ports:
  - "HOST_PORT:CONTAINER_PORT"
```

#### Current Port Configuration
1. **App Service**: `"8080:8080"`
   - Container port: 8080
   - Host port: 8080
   - Access via: `localhost:8080`

2. **PostgreSQL**: `"5432:5432"`
   - Container port: 5432
   - Host port: 5432
   - Access via: `localhost:5432`

3. **Swagger UI**: `"8081:8080"`
   - Container port: 8080
   - Host port: 8081
   - Access via: `localhost:8081`

### Visual Representation
```
Host (Your Computer)        Container
+----------------+         +----------------+
|                |         |                |
|  localhost:8080| <-----> |  container:8080|  (App)
|                |         |                |
|  localhost:5432| <-----> |  container:5432|  (PostgreSQL)
|                |         |                |
|  localhost:8081| <-----> |  container:8080|  (Swagger)
|                |         |                |
+----------------+         +----------------+
```

### Difference Between EXPOSE and Port Mapping

1. **EXPOSE in Dockerfile**
   - Only documents that the container will use the specified port
   - Does not make the port accessible from the host
   - Acts as documentation for container port usage

2. **Port Mapping in docker-compose.yml**
   - Makes container ports accessible from the host
   - Allows changing the host port (e.g., Swagger: 8081:8080)
   - Enables external access to the application

### Benefits of Port Mapping

1. **Isolation**: Containers run in isolation while remaining accessible
2. **Flexibility**: Can change host ports without modifying the application
3. **Security**: Control which ports are exposed to the host
4. **Development**: Facilitates debugging and testing

## Usage

### Development Environment

1. Start all services:
```bash
docker-compose up
```

2. Start specific services:
```bash
docker-compose up app postgres
```

3. Rebuild services:
```bash
docker-compose up --build
```

### Production Environment

1. Build and run using production Dockerfile:
```bash
DOCKERFILE=Dockerfile docker-compose up --build
```

## Environment Variables

The project uses environment variables for configuration. Create `.env.{environment}` files (e.g., `.env.dev`, `.env.sit`, `.env.uat`, `.env.production`) with the following variables:

```env
# Environment (dev, sit, uat, production)
ENVIRONMENT=dev

# Server Configuration
SERVER_PORT=8080

# Database Configuration
# Note: 
# - If running locally (without Docker): use DB_HOST=localhost
# - If running with Docker: use DB_HOST=postgres (service name in docker-compose)
DB_HOST=postgres  # or localhost for local development
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=jeki
```

### Environment File Priority

The application reads environment variables in the following order:
1. First, it attempts to read from `.env.{env}` (e.g., `.env.dev` if `env=dev`)
2. Second, if that file is not found, it tries to read from `.env`
3. Third, if neither file is found, the application will run and read environment variables set by Docker Compose

### Docker Environment Variables

When using Docker, the following environment variables are used with their default values:
- `DB_USER`: postgres (default)
- `DB_PASSWORD`: postgres (default)
- `DB_NAME`: jeki (default)

These variables are read from your `.env.{environment}` file (e.g., `.env.dev`). The same file is used by both the application and PostgreSQL service.

## Volumes

- `postgres_data`: Persistent PostgreSQL data storage
- Application code is mounted as a volume in development mode

## Networks

All services are connected through the `jeki-network` bridge network.

## Health Checks

PostgreSQL service includes health checks to ensure the database is ready before the application starts.

## Best Practices

1. Always use the development setup (`Dockerfile.dev`) for local development
2. Use production setup (`Dockerfile`) for deployment
3. Keep sensitive information in environment files
4. Use the provided health checks to ensure service availability
5. Regularly update base images for security patches

## Troubleshooting

### Database Connection Issues

1. **Connection Refused Error**
   - If running with Docker:
     - Make sure `DB_HOST=postgres` in your `.env.{environment}` file
     - Check if PostgreSQL container is running: `docker-compose ps`
     - Check PostgreSQL logs: `docker-compose logs postgres`
   - If running locally:
     - Make sure `DB_HOST=localhost` in your `.env.{environment}` file
     - Check if PostgreSQL is running:
       ```bash
       # For Mac
       brew services list
       
       # For Ubuntu/Debian
       sudo systemctl status postgresql
       ```

2. **Wrong Database Host**
   - When using Docker: Use `DB_HOST=postgres` (service name in docker-compose)
   - When running locally: Use `DB_HOST=localhost`

3. **Port Conflicts**
   - If you get port conflict errors, make sure:
     - No other PostgreSQL instance is running locally
     - No other application is using port 5432
     - Docker containers are not conflicting with local services

4. **Container Not Starting**
   - Check container logs: `docker-compose logs`
   - Verify environment variables: `docker-compose config`
   - Check for port conflicts: `lsof -i :5432`

5. **Database Initialization Issues**
   - Check initialization logs: `docker-compose logs postgres`
   - Verify init script: `scripts/init.sql`
   - Check permissions on mounted volumes 