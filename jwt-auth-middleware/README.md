# Jeki Backend

Backend service for Jeki application built with Go.

## Table of Contents
- [Jeki Backend](#jeki-backend)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Project Structure](#project-structure)
  - [Getting Started](#getting-started)
    - [1. Environment Setup](#1-environment-setup)
    - [2. Running the Application](#2-running-the-application)
      - [Option 1: Using Docker (Recommended)](#option-1-using-docker-recommended)
      - [Option 2: Running Locally](#option-2-running-locally)
    - [3. Development Workflow](#3-development-workflow)
      - [Using Docker (Recommended)](#using-docker-recommended)
      - [Using Local Setup](#using-local-setup)
  - [Development](#development)
    - [Environment Setup](#environment-setup)
    - [Development Mode (with Hot Reload)](#development-mode-with-hot-reload)
    - [Production Mode](#production-mode)
    - [Manual Mode](#manual-mode)
    - [Environment File Priority](#environment-file-priority)
    - [Available Make Commands](#available-make-commands)
    - [Docker Commands](#docker-commands)
    - [Database Management](#database-management)
      - [Using Docker](#using-docker)
      - [Installing PostgreSQL Client Tools](#installing-postgresql-client-tools)
      - [Accessing Database in Docker](#accessing-database-in-docker)
      - [Local Database Setup](#local-database-setup)
      - [Using psql Command Line](#using-psql-command-line)
      - [Using GUI Tools](#using-gui-tools)
    - [Database Backup and Restore](#database-backup-and-restore)
  - [Future Improvements](#future-improvements)
  - [Troubleshooting](#troubleshooting)
    - [Common Issues and Solutions](#common-issues-and-solutions)
      - [Database Connection Issues](#database-connection-issues)
  - [Development Guidelines](#development-guidelines)
    - [Code Structure](#code-structure)
    - [Git Workflow](#git-workflow)
    - [Testing](#testing)
    - [Code Style](#code-style)
  - [API Documentation](#api-documentation)
    - [Swagger UI](#swagger-ui)
    - [Available Endpoints](#available-endpoints)
    - [API Versioning](#api-versioning)
  - [Environment Variables](#environment-variables)
    - [Required Variables](#required-variables)
    - [Optional Variables](#optional-variables)
    - [Environment File Priority](#environment-file-priority-1)
    - [Docker Environment Variables](#docker-environment-variables)
    - [Authentication Endpoints](#authentication-endpoints)
      - [Google OAuth Login](#google-oauth-login)
      - [Refresh Token](#refresh-token)
      - [Logout](#logout)

## Prerequisites

- Go 1.22 or higher
- Docker and Docker Compose
- PostgreSQL 15
- Make (optional, for using Makefile commands)

## Project Structure

```
.
├── api/            # API documentation and swagger specs
├── bin/            # Binary files
├── cmd/            # Application entry points
├── internal/       # Private application code
│   ├── config/     # Configuration management
│   ├── modules/    # Modular business domains
│   │   ├── user/   # User module
│   │   │   ├── domain/     # Domain models and interfaces
│   │   │   ├── handler/    # HTTP handlers
│   │   │   ├── usecase/    # Business logic
│   │   │   └── repository/ # Data access layer
│   │   ├── auth/   # Authentication module
│   │   │   ├── config/     # Auth configuration
│   │   │   ├── domain/     # Auth models and interfaces
│   │   │   ├── handler/    # Auth HTTP handlers
│   │   │   ├── middleware/ # JWT middleware
│   │   │   ├── repository/ # Session storage
│   │   │   └── usecase/    # Auth business logic
│   │   ├── payment/ # Payment module (future)
│   │   └── notification/ # Notification module (future)
│   ├── pkg/        # Shared libraries
│   │   ├── common/ # Common utilities
│   │   ├── logger/ # Logging utilities
│   │   └── middleware/ # Shared middleware
│   └── handler/    # HTTP routing layer
│       ├── router.go # Entry point for routing (calls v1 router)
│       └── v1/      # API version 1 routes
├── pkg/            # Public library code
├── scripts/        # Build and deployment scripts
├── test/           # Additional test files
├── tmp/            # Temporary files
├── Dockerfile      # Docker configuration
├── docker-compose.yml # Docker compose configuration
├── go.mod          # Go module definition
├── go.sum          # Go module checksums
├── LICENSE         # Project license
├── Makefile        # Build automation
└── README.md       # Project documentation
```

## Getting Started

### 1. Environment Setup

The application supports multiple environments through environment files:

- `.env.dev` - Development environment
- `.env.sit` - System Integration Testing
- `.env.uat` - User Acceptance Testing
- `.env.production` - Production environment

Create a `.env.{environment}` file based on `.env.example`:

```env
# Environment (dev, sit, uat, production)
ENVIRONMENT=dev

# Server Configuration
SERVER_PORT=8080

# Database Configuration
DB_HOST=postgres  # or localhost for local development
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=jeki

# Auth Configuration
GOOGLE_CLIENT_ID=your_client_id
GOOGLE_CLIENT_SECRET=your_client_secret
JWT_SECRET=your_jwt_secret
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=7d
```

**Prioritas Pembacaan File ENV:**

1. **Pertama:** Mencoba membaca file `.env.{env}` (misal: `.env.dev` jika `env=dev`)
2. **Kedua:** Jika file `.env.{env}` tidak ditemukan, mencoba membaca file `.env`
3. **Ketiga:** Jika kedua file tidak ditemukan, aplikasi tetap jalan dan membaca environment variable dari environment (yang di-set oleh Docker Compose)

### 2. Running the Application

You have two options to run the application: using Docker (recommended) or running locally.

#### Option 1: Using Docker (Recommended)

This is the recommended approach as it handles all dependencies automatically, including PostgreSQL. You don't need to install or start PostgreSQL locally.

1. Start the development environment:

   ```bash
   make docker-up
   ```

   This will:

   - Start PostgreSQL in a container
   - Initialize the database automatically
   - Start the application
   - Start Swagger UI

   Note: When using Docker, make sure your `.env.{environment}` file has `DB_HOST=postgres` as this is the service name in docker-compose.

2. Make your code changes - the application will automatically reload

3. Access Swagger UI at <http://localhost:8081> to test your API endpoints

4. When done:

   ```bash
   make docker-down
   ```

Note: If you previously started PostgreSQL locally, you might want to stop it to avoid port conflicts:

```bash
# For Mac
brew services stop postgresql@15

# For Ubuntu/Debian
sudo systemctl stop postgresql
```

#### Option 2: Running Locally

If you prefer to run the application without Docker, follow these steps:

1. Make sure your `.env.{environment}` file has `DB_HOST=localhost` since you're running PostgreSQL locally.

2. Install and start PostgreSQL:

   ```bash
   # For Mac (using Homebrew)
   brew install postgresql@15
   brew services start postgresql@15

   # For Ubuntu/Debian
   sudo apt update
   sudo apt install postgresql-15
   sudo systemctl start postgresql
   ```

3. Initialize the database:

   ```bash
   # First time setup
   make db-setup

   # If you need to reset the database
   make db-reset
   ```

4. Run the application:

   ```bash
   make run
   ```

### 3. Development Workflow

#### Using Docker (Recommended)

1. Start the development environment:

   ```bash
   make docker-up
   ```

2. Make your code changes - the application will automatically reload

3. Access Swagger UI at <http://localhost:8081> to test your API endpoints

4. When done:

   ```bash
   make docker-down
   ```

#### Using Local Setup

1. Start PostgreSQL service:

   ```bash
   # For Mac
   brew services start postgresql@15

   # For Ubuntu/Debian
   sudo systemctl start postgresql
   ```

2. Setup database (This will create the postgres user, database, and run initialization scripts):

   ```bash
   make db-setup
   ```

3. Verify connection:

   ```bash
   psql -U postgres jeki
   ```

Note: The `make db-setup` command will automatically:

- Create the postgres user if it doesn't exist
- Create the jeki database
- Grant necessary permissions
- Run initialization scripts

4. Run the application:

   **Option 1: Standard Run (Manual Restart)**

   ```bash
   make run
   ```

   **Option 2: Hot Reload (Recommended)**

   ```bash
   # Install air for hot reload (using compatible version)
   go install github.com/cosmtrek/air@v1.49.0

   # Add Go bin to PATH if not already added
   echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
   source ~/.zshrc

   # Verify air is installed
   which air

   # Run with hot reload
   air
   ```

   Note: If you get any installation errors, you can also try:

   ```bash
   # Alternative installation method
   curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

   # Add to PATH
   echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
   source ~/.zshrc
   ```

5. Development workflow:

   - Make your code changes
   - With hot reload (Option 2):
     - Changes will be automatically detected
     - Application will restart automatically
   - Without hot reload (Option 1):
     - Press Ctrl+C to stop the current process
     - Run `make run` to restart

6. Access the API at <http://localhost:8080>

Note: If you stop your computer or restart, you'll need to start PostgreSQL service again before running the application. You only need to create the postgres user once.

## Development

### Environment Setup

The application supports multiple environment files:
- `.env.dev` - Development environment
- `.env.sit` - System Integration Testing environment
- `.env.uat` - User Acceptance Testing environment
- `.env.production` - Production environment

Create these files based on the `.env.example` template and configure them according to your environment.

### Development Mode (with Hot Reload)

For development with hot reload (using Air):
```bash
make dev
```
This will:
- Use `Dockerfile.dev`
- Mount source code as volume
- Enable hot reload with Air
- Perfect for development workflow

### Production Mode

For production deployment:
```bash
make prod
```
This will:
- Use `Dockerfile` (multi-stage build)
- Create optimized production image
- No hot reload
- Suitable for production environment

### Manual Mode

You can also specify the Dockerfile manually:
```bash
# Development
DOCKERFILE=Dockerfile.dev docker-compose up --build -d

# Production
DOCKERFILE=Dockerfile docker-compose up --build -d
```

### Environment File Priority

The application reads environment variables in the following order:
1. First, it attempts to read from `.env.{env}` (e.g., `.env.dev` if `env=dev`)
2. Second, if that file is not found, it tries to read from `.env`
3. Third, if neither file is found, the application will run and read environment variables set by Docker Compose

### Available Make Commands

```bash
# Build and run
make build        # Build the application
make run         # Run the application
make docker-build # Build Docker image
make docker-up   # Start Docker containers
make docker-down # Stop Docker containers (data tetap ada)
make docker-down-volumes # Stop containers dan hapus semua data

# Database
make db-setup    # Setup database (untuk local setup)
make db-setup-docker # Setup database (untuk Docker setup)
make db-reset    # Reset database (untuk local setup)
make db-reset-docker # Reset database (untuk Docker setup)

# Development tools
make test        # Run tests
make swagger     # Generate Swagger docs
make deps        # Install dependencies
make clean       # Clean build artifacts
```

### Docker Commands

```bash
# Build and start containers
make docker-up

# Stop and remove containers
make docker-down

# View application logs
make logs

# Rebuild and restart containers
make docker-build
```

### Database Management

#### Using Docker
Before running database commands with Docker, make sure to start the containers first:

```bash
# Start Docker containers
make docker-up

# Then you can run database commands
make db-setup-docker  # Setup database
make db-reset-docker  # Reset database
```

Note: Always run `make docker-up` before running any database commands that use Docker (`db-setup-docker` or `db-reset-docker`). This ensures that the PostgreSQL container is running and ready to accept connections.

#### Installing PostgreSQL Client Tools

Before using psql or other PostgreSQL client tools, you need to install them:

```bash
# For Mac (using Homebrew)
brew install postgresql@15

# For Ubuntu/Debian
sudo apt update
sudo apt install postgresql-client-15
```

After installation, you might need to add PostgreSQL to your PATH:

```bash
# For Mac
echo 'export PATH="/opt/homebrew/opt/postgresql@15/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# For Ubuntu/Debian
echo 'export PATH="/usr/lib/postgresql/15/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

#### Accessing Database in Docker

When using Docker, you can access the database in several ways:

1. **Using psql inside Docker container**:
   ```bash
   # Connect directly to jeki database
   docker-compose exec postgres psql -U postgres jeki

   # Or connect to postgres first
   docker-compose exec postgres psql -U postgres
   ```

2. **Using GUI Tools**:
   - **pgAdmin**:
     - Host: localhost
     - Port: 5432
     - Database: jeki
     - Username: postgres
     - Password: postgres

   - **DBeaver**:
     - Host: localhost
     - Port: 5432
     - Database: jeki
     - Username: postgres
     - Password: postgres

3. **Common psql commands** (when connected):
   ```sql
   \l          # List databases
   \c jeki     # Connect to jeki database
   \dt         # List tables
   \d users    # Describe users table
   \q          # Quit psql
   ```

4. **Database Management Commands**:
   ```bash
   # Setup database
   make db-setup-docker

   # Reset database
   make db-reset-docker

   # Backup database
   docker-compose exec postgres pg_dump -U postgres jeki > backup.sql

   # Restore database
   docker-compose exec -T postgres psql -U postgres jeki < backup.sql
   ```

#### Local Database Setup

If you're running PostgreSQL locally (without Docker), follow these steps in order:

1. **Start PostgreSQL**:

   ```bash
   # For Mac
   brew services start postgresql@15

   # For Ubuntu/Debian
   sudo systemctl start postgresql
   ```

2. **Setup database** (This will create the postgres user, database, and run initialization scripts):

   ```bash
   make db-setup
   ```

3. **Verify connection**:

   ```bash
   psql -U postgres jeki
   ```

Note: The `make db-setup` command will automatically:

- Create the postgres user if it doesn't exist
- Create the jeki database
- Grant necessary permissions
- Run initialization scripts

#### Using psql Command Line

```bash
# Connect to database
# psql -U <username> <database_name>
psql -U postgres jeki

# Common psql commands
\l          # List databases
\c jeki     # Connect to jeki database
\dt         # List tables
\d users    # Describe users table
\q          # Quit psql (exit to terminal)
```

Note: Commands starting with `\` are psql internal commands, not shell commands. For example, `\q` is used to exit psql and return to your terminal.

#### Using GUI Tools

1. **pgAdmin** (Official PostgreSQL GUI)

   - Download from: <https://www.pgadmin.org/download/>
   - Connect using:
     - Host: localhost
     - Port: 5432
     - Database: jeki
     - Username: postgres
     - Password: postgres

2. **DBeaver** (Multi-database GUI)
   - Download from: <https://dbeaver.io/>
   - Support for multiple databases
   - User-friendly interface

### Database Backup and Restore

```bash
# Backup
pg_dump -U postgres jeki > backup.sql

# Restore
psql -U postgres jeki < backup.sql
```

## Future Improvements

Berikut adalah saran pengembangan untuk masa depan:

1. **Testing**:
   - Tambahkan folder `test` di setiap modul
   - Implementasikan unit test, integration test, dan e2e test
   - Gunakan test coverage tools

2. **API Documentation**:
   - Tambahkan Swagger/OpenAPI documentation
   - Dokumentasikan setiap endpoint
   - Tambahkan contoh request/response

3. **Error Handling**:
   - Buat package `errors` di `pkg/`
   - Standardisasi error response
   - Implementasikan error logging

4. **Logging**:
   - Buat package `logger` di `pkg/`
   - Standardisasi logging format
   - Implementasikan log rotation

5. **Monitoring**:
   - Tambahkan health check endpoint
   - Implementasikan metrics (Prometheus)
   - Setup tracing (OpenTelemetry)

6. **Security**:
   - Implementasikan rate limiting
   - Setup CORS
   - Implementasikan security headers

7. **CI/CD**:
   - Setup GitHub Actions/GitLab CI
   - Implementasikan automated testing
   - Setup automated deployment

8. **Containerization**:
   - Optimize Dockerfile
   - Setup multi-stage build
   - Implementasikan health check

9. **Database**:
   - Implementasikan migrations
   - Setup database backup
   - Implementasikan connection pooling

10. **Caching**:
    - Implementasikan Redis
    - Setup cache invalidation
    - Implementasikan cache warming

11. **Message Queue**:
    - Setup RabbitMQ/Kafka
    - Implementasikan async processing
    - Setup dead letter queue

12. **Service Discovery**:
    - Setup service registry
    - Implementasikan load balancing
    - Setup circuit breaker

13. **Configuration Management**:
    - Implementasikan feature flags
    - Setup configuration validation
    - Implementasikan secrets management

14. **Monitoring & Alerting**:
    - Setup Grafana dashboards
    - Implementasikan alerting
    - Setup log aggregation

15. **Documentation**:
    - Tambahkan API documentation
    - Dokumentasikan deployment process
    - Tambahkan troubleshooting guide
## Troubleshooting

### Common Issues and Solutions

#### Database Connection Issues

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

## Development Guidelines

### Code Structure
- Follow the modular monolith architecture
- Keep business logic in usecase layer
- Use interfaces for better testability
- Implement proper error handling

### Git Workflow
1. Create feature branch from main
2. Make changes and commit
3. Write tests if applicable
4. Create pull request
5. Get code review
6. Merge to main

### Testing
- Write unit tests for business logic
- Use table-driven tests
- Mock external dependencies
- Aim for good test coverage

### Code Style
- Use `gofmt` for formatting
- Follow Go best practices
- Write clear comments
- Use meaningful variable names

## API Documentation

### Swagger UI
Access the API documentation at:
- Development: http://localhost:8081
- Production: https://api.jeki.com/docs

### Available Endpoints
1. **Health Check**:
   - GET `/ping`
   - Response: `{"message": "pong"}`

2. **User Management**:
   - POST `/api/v1/users` - Create user
   - GET `/api/v1/users` - List users
   - GET `/api/v1/users/:id` - Get user details
   - PUT `/api/v1/users/:id` - Update user
   - DELETE `/api/v1/users/:id` - Delete user

### API Versioning
- Current version: v1
- Version prefix: `/api/v1/`
- Future versions will use `/api/v2/`, etc.

## Environment Variables

### Required Variables
```env
# Environment (dev, sit, uat, production)
ENVIRONMENT=dev

# Server Configuration
SERVER_PORT=8080

# Database Configuration
DB_HOST=postgres  # or localhost for local development
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=jeki

# PostgreSQL Docker Configuration (optional, has defaults)
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=jeki
```

### Optional Variables
```env
# Logging
LOG_LEVEL=debug  # debug, info, warn, error
LOG_FORMAT=json  # json, text

# Security
JWT_SECRET=your_jwt_secret
JWT_EXPIRY=24h

# Rate Limiting
RATE_LIMIT=100
RATE_WINDOW=1m
```

### Environment File Priority
1. `.env.{environment}` (e.g., `.env.dev`)
2. `.env`
3. System environment variables

### Docker Environment Variables
When using Docker, the following environment variables are used with their default values:
- `DB_USER`: postgres (default)
- `DB_PASSWORD`: postgres (default)
- `DB_NAME`: jeki (default)

These variables are read from your `.env.{environment}` file (e.g., `.env.dev`). The same file is used by both the application and PostgreSQL service.

Example `.env.dev`:
```env
# Environment (dev, sit, uat, production)
ENVIRONMENT=dev

# Server Configuration
SERVER_PORT=8080

# Database Configuration
DB_HOST=postgres  # or localhost for local development
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=jeki
```

Note: If these variables are not set in your `.env` file, the default values will be used:
- DB_USER=postgres
- DB_PASSWORD=postgres
- DB_NAME=jeki

### Authentication Endpoints

#### Google OAuth Login
```http
POST /v1/auth/google?code={authorization_code}
```

#### Refresh Token
```http
POST /v1/auth/refresh?refresh_token={refresh_token}
```

#### Logout
```http
POST /v1/auth/logout
Authorization: Bearer {access_token}
```

For detailed API documentation, see:
- [Auth Module Documentation](internal/modules/auth/README.md)
- [Swagger UI](http://localhost:8081) (when running in development mode)

