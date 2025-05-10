package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	// Get database configuration from environment variables with defaults
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5433")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "postgres")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	fmt.Println("Connected to Database...")
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	// First, drop the goose version table if it exists
	_, err := db.Exec("DROP TABLE IF EXISTS goose_db_version")
	if err != nil {
		return fmt.Errorf("failed to drop goose version table: %w", err)
	}

	// Drop all existing tables
	_, err = db.Exec(`
		DROP TABLE IF EXISTS users CASCADE;
		DROP TABLE IF EXISTS workouts CASCADE;
		DROP TABLE IF EXISTS workout_entries CASCADE;
	`)
	if err != nil {
		return fmt.Errorf("failed to drop existing tables: %w", err)
	}

	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	// Set the base filesystem for goose
	goose.SetBaseFS(migrationsFS)

	// Run the migrations
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	// Reset the base filesystem
	goose.SetBaseFS(nil)
	return nil
}

func Migrate(db *sql.DB, dir string) error {
	// First, drop the goose version table if it exists
	_, err := db.Exec("DROP TABLE IF EXISTS goose_db_version")
	if err != nil {
		return fmt.Errorf("failed to drop goose version table: %w", err)
	}

	// Drop all existing tables
	_, err = db.Exec(`
		DROP TABLE IF EXISTS users CASCADE;
		DROP TABLE IF EXISTS workouts CASCADE;
		DROP TABLE IF EXISTS workout_entries CASCADE;
	`)
	if err != nil {
		return fmt.Errorf("failed to drop existing tables: %w", err)
	}

	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}
