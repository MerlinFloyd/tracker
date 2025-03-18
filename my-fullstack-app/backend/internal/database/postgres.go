package database

import (
	"database/sql"
	"fmt"
	"log"

	"my-fullstack-app/backend/internal/models"

	_ "github.com/lib/pq"
)

const (
	host     = "db"
	port     = 5432
	user     = "app"
	password = "yourpassword"
	dbname   = "appdb"
)

func Connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error opening database: ", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
		return nil, err
	}

	log.Println("Successfully connected to the  database")
	return db, nil
}

// StoreBalance stores an Ethereum balance record in the database
func StoreBalance(db *sql.DB, record models.BalanceRecord) (int, error) {
	// SQL query to insert a balance record
	query := `
        INSERT INTO balance_records (address, balance, balance_eth, fetched_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	// Execute the query
	var id int
	err := db.QueryRow(
		query,
		record.Address,
		record.Balance,
		record.BalanceETH,
		record.FetchedAt,
	).Scan(&id)

	if err != nil {
		log.Printf("Error storing balance record: %v", err)
		return 0, err
	}

	return id, nil
}
