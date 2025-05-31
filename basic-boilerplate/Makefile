.PHONY: build run test clean docker-build docker-up docker-down dev prod go-mod-tidy db-setup db-setup-docker db-reset db-reset-docker swagger deps

# Build the application
build:
	go build -o bin/api ./cmd/api

# Run the application
run:
	@echo "Checking if port 8080 is in use..."
	@if lsof -i :8080 > /dev/null; then \
		echo "Port 8080 is in use. Attempting to free it..."; \
		lsof -ti :8080 | xargs kill -9 2>/dev/null || true; \
		echo "Port 8080 has been freed."; \
	fi
	go run ./cmd/api

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Docker commands
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Development environment (with hot reload)
dev:
	DOCKERFILE=Dockerfile.dev docker-compose up --build -d

# Production environment
prod:
	DOCKERFILE=Dockerfile docker-compose up --build -d

# Setup database (without Docker)
db-setup:
	@echo "Setting up database..."
	@echo "Checking PostgreSQL service..."
	@if [ "$(shell uname)" = "Darwin" ]; then \
		brew services list | grep postgresql@15 > /dev/null || (echo "PostgreSQL is not running. Please start it first with: brew services start postgresql@15" && exit 1); \
	else \
		systemctl is-active --quiet postgresql || (echo "PostgreSQL is not running. Please start it first with: sudo systemctl start postgresql" && exit 1); \
	fi
	@echo "Creating postgres user if not exists..."
	@psql -c "CREATE USER postgres WITH PASSWORD 'postgres' SUPERUSER;" postgres || true
	@echo "Creating database..."
	@psql -U postgres -c "CREATE DATABASE jeki;" || true
	@echo "Granting privileges..."
	@psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE jeki TO postgres;"
	@echo "Running initialization scripts..."
	@psql -U postgres jeki -f scripts/init.sql
	@echo "Database setup completed successfully!"
	@echo "\nTo connect to the database, use:"
	@echo "psql -U postgres jeki"
	@echo "\nOr to connect as postgres user first:"
	@echo "psql postgres"
	@echo "Then in psql, you can use:"
	@echo "\\c jeki    # to connect to jeki database"
	@echo "\\dt        # to list tables"
	@echo "\\q         # to quit psql"

# Reset database (without Docker)
db-reset:
	@echo "Resetting database..."
	@echo "Terminating all connections to database..."
	@psql -U postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'jeki';" || true
	@echo "Dropping database..."
	@psql -U postgres -c "DROP DATABASE IF EXISTS jeki;"
	@echo "Setting up fresh database..."
	@make db-setup

# Reset database (with Docker)
db-reset-docker:
	@echo "Resetting database in Docker..."
	@docker-compose exec postgres psql -U postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'jeki';" || true
	@docker-compose exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS jeki;"
	@docker-compose exec postgres psql -U postgres -c "CREATE DATABASE jeki;"
	@docker-compose exec postgres psql -U postgres jeki -f /docker-entrypoint-initdb.d/init.sql
	@echo "Database has been reset successfully!"

# Generate Swagger docs
swagger:
	swag init -g cmd/api/main.go -o docs/swagger

# Install dependencies
deps:
	go mod download

# Tidy up Go module dependencies
go-mod-tidy:
	@echo "Tidying up Go module dependencies..."
	@go mod tidy
	@echo "Dependencies have been tidied up successfully!"

# Setup database (with Docker)
db-setup-docker:
	@echo "Setting up database in Docker..."
	@docker-compose exec postgres psql -U postgres -c "CREATE DATABASE jeki;" || true
	@docker-compose exec postgres psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE jeki TO postgres;"
	@docker-compose exec postgres psql -U postgres jeki -f /docker-entrypoint-initdb.d/init.sql
	@echo "Database setup completed successfully!"
	@echo "\nTo connect to the database in Docker, use:"
	@echo "docker-compose exec postgres psql -U postgres jeki"
	@echo "\nOr to connect to psql first:"
	@echo "docker-compose exec postgres psql -U postgres"
	@echo "Then in psql, you can use:"
	@echo "\\c jeki    # to connect to jeki database"
	@echo "\\dt        # to list tables"
	@echo "\\q         # to quit psql"
