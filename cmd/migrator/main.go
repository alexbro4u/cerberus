package main

import (
	"cerberus/internal/config"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
)

func main() {
	cfg := config.MustLoadMigrator()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s?x-migrations-table=%s&sslmode=disable",
		cfg.Storage.User, cfg.Storage.Password, cfg.Storage.Host, cfg.Storage.DbName, cfg.MigrationsTable)

	migrationsPath := fmt.Sprintf("file://%s", cfg.MigrationsPath)

	m, err := migrate.New(
		migrationsPath,
		dbURL,
	)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No changes to migrate")
			return
		}
		panic(err)
	}

	fmt.Println("Migrations applied")
}
