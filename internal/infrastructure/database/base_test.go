package database

import (
	"os"
	"testing"

	"chatapp/internal/domain/entity"

	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	// Setup test database
	dbConfig, err := NewPostgresConfig(
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_SSLMODE"),
	)
	if err != nil {
		panic(err)
	}

	testDSN := dbConfig.NewDSN()
	testDB, err = NewPostgresDB(testDSN)
	if err != nil {
		panic(err)
	}

	// Migrate test database
	testDB.AutoMigrate(&entity.User{})

	// Tear down test database
	defer func() {
		if err := testDB.Migrator().DropTable(&entity.User{}); err != nil {
			panic(err)
		}
	}()

	// Run tests
	exitCode := m.Run()
	os.Exit(exitCode)
}
