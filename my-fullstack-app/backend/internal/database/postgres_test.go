package database

import (
	"os"
	"testing"
)

func TestConnect(t *testing.T) {
	// Skip if not in test environment to avoid connecting to production DB
	if os.Getenv("GO_ENV") != "test" {
		t.Skip("Skipping database connection test in non-test environment")
	}

	// Attempt to connect to the database
	db, err := Connect()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test that we can ping the database
	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}
