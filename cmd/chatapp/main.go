package main

import (
	"log"
	"os"

	"chatapp/internal/domain/entity"
	"chatapp/internal/infrastructure/database"
	"chatapp/internal/interface/router"

	"github.com/labstack/echo/v4"
)

func main() {
	// Create a new database connection
	dbConfig, err := database.NewPostgresConfig(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
	if err != nil {
		log.Fatal(err)
	}

	dsn := dbConfig.NewDSN()
	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to database:", db.Name())

	// Migrate the database
	database.Migrate(db, &entity.User{})
	log.Println("Successfully migrated database")

	// Set up router
	e := echo.New()
	handlers := router.InitRouter(db)
	handlers.SetUpRouter(e)

	e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))
}
