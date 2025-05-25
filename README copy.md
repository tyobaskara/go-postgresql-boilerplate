# Jeki Backend

Backend service for Jeki application built with Go.

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
│   ├── domain/     # Domain models and interfaces
│   ├── handler/    # HTTP handlers
│   ├── usecase/    # Business logic
│   ├── infrastructure/ # Infrastructure implementations (database, cache, etc)
│   └── delivery/   # Delivery layer (API implementations)
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
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=jeki
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

1. Install and start PostgreSQL:

   ```bash
   # For Mac (using Homebrew)
   brew install postgresql@15
   brew services start postgresql@15

   # For Ubuntu/Debian
   sudo apt update
   sudo apt install postgresql-15
   sudo systemctl start postgresql
   ```

2. Initialize the database:

   ```bash
   # First time setup
   make db-setup

   # If you need to reset the database
   make db-reset
   ```

3. Run the application:

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
make db-setup    # Setup database
make db-reset    # Reset database

# Development tools
make test        # Run tests
make swagger     # Generate Swagger docs
make deps        # Install dependencies
make clean       # Clean build artifacts
```

### Docker Commands

The application uses Docker for containerization. Here are the key Docker-related commands:

1. **Starting Services**:

   ```bash
   make docker-up
   ```

   This will:

   - Build the application
   - Start PostgreSQL database
   - Initialize the database with required schemas
   - Start the application
   - Start Swagger UI

2. **Stopping Services**:

   ```bash
   # Stop services but keep the data
   make docker-down

   # Stop services and remove all data (including database)
   make docker-down-volumes
   ```

3. **Accessing Services**:
   - API: <http://localhost:8080>
   - Swagger UI: <http://localhost:8081>
   - Database: localhost:5432

Note: Using `make docker-down` is safe for development as it preserves your database data. Use `make docker-down-volumes` only when you want to completely reset your environment.

### Database Management

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