package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations applies all migrations in the specified directory
func RunMigrations(db *sql.DB, migrationsDir string) error {
	// Create migrations table if it doesn't exist
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS migrations (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of applied migrations
	appliedMigrations := make(map[string]bool)
	rows, err := db.Query("SELECT name FROM migrations")
	if err != nil {
		return fmt.Errorf("failed to query migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan migration row: %w", err)
		}
		appliedMigrations[name] = true
	}

	// Read migration files
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter and sort SQL files
	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	// Apply each migration
	for _, filename := range migrationFiles {
		if appliedMigrations[filename] {
			log.Printf("Migration %s already applied, skipping", filename)
			continue
		}

		log.Printf("Applying migration: %s", filename)

		// Read migration file
		filePath := filepath.Join(migrationsDir, filename)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Begin transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		// Execute the migration
		_, err = tx.Exec(string(content))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		// Record the migration
		_, err = tx.Exec("INSERT INTO migrations (name) VALUES ($1)", filename)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", filename, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		log.Printf("Successfully applied migration: %s", filename)
	}

	return nil
}
