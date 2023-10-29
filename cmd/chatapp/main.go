package main

import (
	"log"
	"os"

	"chatapp/internal/infrastructure/database"
	"chatapp/internal/interface/handler"
	"chatapp/internal/interface/router"
	"chatapp/internal/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	// Create a new database connection
	dbConfig := database.NewPostgresConfig(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL_MODE"),
	)
	if err := dbConfig.Validate(); err != nil {
		log.Fatal(err)
	}

	dsn := dbConfig.NewDSN()
	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to database:", db.Name())

	// Migrate the database
	database.Migrate(db)
	log.Println("Successfully migrated database")

	// Set up router
	e := echo.New()

	userRepo := database.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	authHandler := handler.NewAuthHandler(userUseCase)
	router := router.NewRouter(*authHandler)
	router.NewRouter(e)

	e.Logger.Fatal(e.Start(":3000"))
}
