package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	// Try to connect to the database
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Printf("Warning: cannot open database connection: %v", err)
		log.Printf("Database tests will be skipped. Make sure PostgreSQL is running and database 'simple_bank' exists.")
		log.Printf("You can start the database with: make up")
		log.Printf("You can create the database with: make createdb")
		log.Printf("You can run migrations with: make migrateup")

		// Set testDB to nil so tests can check if it's available
		testDB = nil
		testQueries = nil

		// Exit with success since this is expected in CI without database
		os.Exit(0)
	}

	// Test the connection
	err = testDB.Ping()
	if err != nil {
		log.Printf("Warning: cannot ping database: %v", err)
		log.Printf("Database tests will be skipped. Make sure PostgreSQL is running and database 'simple_bank' exists.")

		testDB.Close()
		testDB = nil
		testQueries = nil

		// Exit with success since this is expected in CI without database
		os.Exit(0)
	}

	testQueries = New(testDB)

	// Run the tests
	exitCode := m.Run()

	// Clean up
	if testDB != nil {
		testDB.Close()
	}

	os.Exit(exitCode)
}

// Helper function to skip tests if database is not available
func requireDB(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available, skipping test")
	}
}
