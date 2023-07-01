package main

import (
	"fmt"
	"log"
	"os"
	"time"

	// "github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type TodoStatus string

const (
	TodoStatusPending   TodoStatus = "pending"
	TodoStatusCompleted TodoStatus = "completed"
)

type Todo struct {
	Id          string     `json:"id"`
	Todo        string     `json:"todo"`
	Status      TodoStatus `json:"status"`
	DateCreated time.Time  `json:"date_created"`
}

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load DB variables from .env	
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
					dbUsername, dbPassword, dbHost, dbPort, dbName)
	
	// TODO Connect to DB using connection string

	fmt.Printf("Database URL: %s\n", dbUrl)
	os.Exit(0)
}
