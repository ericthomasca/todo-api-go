package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func connectToDatabase() (*pgxpool.Pool, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	pgURI := os.Getenv("PG_URI")
	if pgURI == "" {
		return nil, fmt.Errorf("PG_URI environment variable is not set: %w", err)
	}

	pgxConfig, err := pgxpool.ParseConfig(pgURI)
	if err != nil {
		return nil, fmt.Errorf("problem parsing PG_URI: %w", err)
	}

	pgConnPool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return pgConnPool, nil
}