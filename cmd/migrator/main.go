package main

import (
	"cerberus/internal/config"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.MustLoadMigrator()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable&x-migrations-table=%s",
		cfg.Storage.User,
		cfg.Storage.Password,
		cfg.Storage.Host,
		cfg.Storage.DbName,
		cfg.MigrationsTable,
	)

	migrationsPath := fmt.Sprintf("file://%s", cfg.MigrationsPath)

	m, err := migrate.New(
		migrationsPath,
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v\n", err)
	}

	log.Println("Starting migration...")
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No changes to migrate")
			return
		}
		log.Fatalf("Migration failed: %v\n", err)
	}

	log.Println("Migrations applied successfully")
}
