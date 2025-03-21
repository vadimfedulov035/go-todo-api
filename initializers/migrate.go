package initializers

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/stdlib"
)

func runMigrations() error {
	if Pool == nil {
		return fmt.Errorf("Database connection pool not initialized")
	}

	// Create stdlib-compatible DB connection
	db := stdlib.OpenDBFromPool(Pool)
	defer db.Close()

	// Create migration driver with custom config
	driver, err := postgres.WithInstance(db, &postgres.Config{
		SchemaName:            "private",
		MigrationsTable:       "schema_migrations",
		MultiStatementEnabled: true,
	})
	if err != nil {
		return fmt.Errorf("Failed to create migration driver: %w", err)
	}

	// Initialize migrator
	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("Failed to initialize migrator: %w", err)
	}
	defer m.Close()

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Migration failed: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}
