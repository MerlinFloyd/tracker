package database

import (
	"database/sql"
	"log"
	"my-fullstack-app/backend/internal/models"
	"time"
)

// StoreTokenBalance stores an ERC20 token balance record in the database
func StoreTokenBalance(db *sql.DB, record models.TokenBalanceRecord) (int, error) {
	// SQL query to insert a token balance record
	query := `
        INSERT INTO balance_records (
            address, balance, balance_eth, fetched_at
        )
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
		log.Printf("Error storing token balance record: %v", err)
		return 0, err
	}

	return id, nil
}

// GetLatestTokenBalances retrieves latest token balances for an address
func GetLatestTokenBalances(db *sql.DB, address string) ([]models.TokenBalanceRecord, error) {
	query := `
        SELECT DISTINCT ON (token_address) 
            id, address, balance, balance_eth, fetched_at
        FROM balance_records
        WHERE address = $1
        ORDER BY address, fetched_at DESC
    `

	rows, err := db.Query(query, address)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.TokenBalanceRecord

	for rows.Next() {
		var record models.TokenBalanceRecord
		err := rows.Scan(
			&record.ID,
			&record.Address,
			&record.Balance,
			&record.BalanceETH,
			&record.FetchedAt,
		)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

// GetTokenBalances retrieves all token balance records for a specific token address
func GetTokenBalances(db *sql.DB, tokenAddress string) ([]models.TokenBalanceRecord, error) {
	// SQL query to retrieve token balances
	query := `SELECT id, address, balance, balance_eth, fetched_at 
	FROM balance_records 
	WHERE address = $1
    `

	// Execute the query
	rows, err := db.Query(query, tokenAddress)
	if err != nil {
		log.Printf("Error retrieving token balances: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Parse results
	var balances []models.TokenBalanceRecord
	for rows.Next() {
		var balance models.TokenBalanceRecord
		var fetchedAt time.Time

		err := rows.Scan(
			&balance.ID,
			&balance.Address,
			&balance.Balance,
			&balance.BalanceETH,
			&fetchedAt,
		)
		if err != nil {
			log.Printf("Error scanning token balance record: %v", err)
			continue
		}

		balance.FetchedAt = fetchedAt
		balances = append(balances, balance)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating token balance records: %v", err)
		return nil, err
	}

	return balances, nil
}
