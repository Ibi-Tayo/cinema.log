# Database Migrations with Goose

This project now uses [Goose](https://github.com/pressly/goose) for database migrations. Goose is a database migration tool written in Go that supports multiple databases including PostgreSQL.

## Migration Files

Migration files are located in `internal/migration/goose/` and follow the naming convention:

```
YYYYMMDDHHMMSS_migration_name.sql
```

Each migration file contains both `Up` and `Down` migrations:

```sql
-- +goose Up
CREATE TABLE example (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- +goose Down
DROP TABLE example;
```

## Migration Commands

### Run All Pending Migrations

```bash
make migrate-up
```

### Rollback Last Migration

```bash
make migrate-down
```

### Check Migration Status

```bash
make migrate-status
```

### Reset Database (Rollback All Migrations)

```bash
make migrate-reset
```

### Create New Migration

```bash
make migrate-create name=add_users_email_column
```

## Manual Migration Management

You can also run migrations manually using the CLI tool:

```bash
# Run all pending migrations
go run cmd/migrate/main.go -command=up

# Rollback last migration
go run cmd/migrate/main.go -command=down

# Check status
go run cmd/migrate/main.go -command=status

# Reset database
go run cmd/migrate/main.go -command=reset
```

## Embedded Migrations

Migrations are embedded into the binary using Go's `embed` package, so you don't need to distribute migration files separately from your application.

## Automatic Migrations on Server Start

The server automatically runs pending migrations when it starts up. If you prefer manual control, you can use the `database.New()` function instead of `database.NewWithMigrations()`.

## Migration History

Goose tracks migration history in a `goose_db_version` table in your database. This table records which migrations have been applied and when.

## Best Practices

1. **Always test migrations**: Test both up and down migrations in a development environment
2. **Keep migrations atomic**: Each migration should be a single, atomic operation
3. **Don't modify existing migrations**: Once a migration is committed and deployed, create a new migration to make changes
4. **Use transactions**: Goose automatically wraps each migration in a transaction (for supported databases like PostgreSQL)
5. **Backup before production migrations**: Always backup your production database before running migrations

## Converting from Old Migration System

Your existing migrations have been converted to Goose format:

- `001_initial_users_table.up.sql` and `001_initial_users_table.down.sql` → `20240817000001_initial_users_table.sql`

The old migration files can be safely removed after confirming the new system works correctly.
