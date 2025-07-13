# Migration Workflow

## Kapan Menggunakan Migration Files

### ✅ Gunakan Migration Files untuk:
- **Production Environment** - Database sudah ada data penting
- **Development yang sudah berjalan** - Ada data yang tidak boleh hilang
- **Team Development** - Perubahan schema harus bisa di-track
- **Database yang sudah live** - Tidak bisa drop dan recreate

### ❌ Jangan Gunakan Migration Files untuk:
- **First-time setup** - Bisa pakai `scripts/init.sql`
- **Development awal** - Bisa pakai `scripts/init.sql`
- **Testing environment** - Bisa pakai `scripts/init.sql`

## Step-by-Step Migration Workflow

### 1. **Setup Environment**

```bash
# Start containers
make docker-up

# Wait for database to be ready
sleep 10
```

### 2. **Create New Migration**

```bash
# Create new migration file
migrate create -ext sql -dir migrations -seq add_new_table

# This will create:
# - migrations/000003_add_new_table.up.sql
# - migrations/000003_add_new_table.down.sql
```

### 3. **Write Migration Files**

**Up Migration** (`000003_add_new_table.up.sql`):
```sql
-- Add new table
CREATE TABLE IF NOT EXISTS new_table (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes
CREATE INDEX idx_new_table_name ON new_table(name);
```

**Down Migration** (`000003_add_new_table.down.sql`):
```sql
-- Drop new table
DROP TABLE IF EXISTS new_table CASCADE;
```

### 4. **Run Migrations**

#### Local Development:
```bash
# Run all pending migrations
make migrate-up-local

# Rollback last migration
make migrate-down-local

# Rollback all migrations
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/maxwash?sslmode=disable" down -all
```

#### Docker Environment:
```bash
# Run all pending migrations
make migrate-up-docker

# Rollback last migration
make migrate-down-docker

# Rollback all migrations
docker run --network maxwash-backend_maxwash-network \
  -v $(pwd)/migrations:/migrations \
  migrate/migrate \
  -path /migrations \
  -database "postgres://postgres:postgres@postgres:5432/maxwash?sslmode=disable" down -all
```

### 5. **Check Migration Status**

```bash
# Check migration status
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/maxwash?sslmode=disable" version

# Check migration status in Docker
docker run --network maxwash-backend_maxwash-network \
  -v $(pwd)/migrations:/migrations \
  migrate/migrate \
  -path /migrations \
  -database "postgres://postgres:postgres@postgres:5432/maxwash?sslmode=disable" version
```

## Best Practices

### 1. **Migration File Naming**
- Use descriptive names: `add_user_roles`, `create_payment_table`
- Use sequential numbers: `000001_`, `000002_`, etc.

### 2. **Migration Content**
- Always include both UP and DOWN migrations
- Use `IF NOT EXISTS` for CREATE statements
- Use `CASCADE` for DROP statements
- Include proper indexes
- Add comments for complex migrations

### 3. **Testing Migrations**
```bash
# Test up migration
make migrate-up-docker

# Test down migration
make migrate-down-docker

# Test up migration again
make migrate-up-docker
```

### 4. **Production Deployment**
```bash
# 1. Backup database first
pg_dump -U postgres maxwash > backup_$(date +%Y%m%d_%H%M%S).sql

# 2. Run migrations
make migrate-up-docker

# 3. Verify application works
curl http://localhost:8080/health

# 4. If issues, rollback
make migrate-down-docker
```

## Common Commands

### Development Workflow:
```bash
# Start fresh with migrations
make docker-up-with-migrations

# Add new migration
migrate create -ext sql -dir migrations -seq add_feature_x

# Edit migration files
# Then run:
make migrate-up-docker

# Check status
migrate -path migrations -database "postgres://postgres:postgres@postgres:5432/maxwash?sslmode=disable" version
```

### Production Workflow:
```bash
# 1. Backup
pg_dump -U postgres maxwash > backup.sql

# 2. Deploy new code
git pull origin main

# 3. Run migrations
make migrate-up-docker

# 4. Verify
curl http://localhost:8080/health

# 5. If rollback needed
make migrate-down-docker
```

## Troubleshooting

### Migration Failed:
```bash
# Check migration status
migrate -path migrations -database "postgres://postgres:postgres@postgres:5432/maxwash?sslmode=disable" version

# Force version (if needed)
migrate -path migrations -database "postgres://postgres:postgres@postgres:5432/maxwash?sslmode=disable" force VERSION

# Check logs
docker-compose logs postgres
```

### Database Connection Issues:
```bash
# Check if database is running
docker-compose ps

# Restart database
docker-compose restart postgres

# Check database logs
docker-compose logs postgres
``` 