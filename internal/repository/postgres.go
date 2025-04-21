package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// NewPostgresDB creates a new connection to PostgreSQL
func NewPostgresDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging the database: %v", err)
	}

	return db, nil
}

type Repositories struct {
	Users     UserRepository
	Accounts  AccountRepository
	Cards     CardRepository
	Credits   CreditRepository
	Transfers TransferRepository
	Analytics AnalyticsRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users:     NewUserRepository(db),
		Accounts:  NewAccountRepository(db),
		Cards:     NewCardRepository(db),
		Credits:   NewCreditRepository(db),
		Transfers: NewTransferRepository(db),
		Analytics: NewAnalyticsRepository(db),
	}
}
