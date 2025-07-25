package postgresdb

import "github.com/jackc/pgx/v5/pgxpool"

type PostgresDB struct {
	db *pgxpool.Pool
}

func New() (*PostgresDB, error) {
	// Initialize the PostgresDB connection here
	// Return the instance of PostgresDB
	return &PostgresDB{}, nil
}
