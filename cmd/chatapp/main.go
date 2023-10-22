package main

import (
	"fmt"
	"log"

	"chatapp/internal/infrastructure/db"
)

func main() {
	// Create a new database connection
	dbConfig := db.NewPostgresConfig()
	if err := dbConfig.Validate(); err != nil {
		log.Fatal(err)
	}

	dsn := dbConfig.NewDSN()
	db, err := db.NewPostgresDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to database:", db.Name())
}
