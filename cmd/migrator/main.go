package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/approachi/grinex-rates/internal/config"
	"github.com/approachi/grinex-rates/internal/migrations"
)

func main() {
	command := parseFlags()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := openDB(cfg)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db, command); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}

func parseFlags() string {
	command := flag.String("command", "up", "goose command: up, down, status, version")
	flag.Parse()
	return *command
}

func openDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
		cfg.Postgres.SSLMode,
	)
	return sql.Open("postgres", dsn)
}

func runMigrations(db *sql.DB, command string) error {
	goose.SetBaseFS(migrations.Files)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.RunContext(nil, command, db, ".")
}
