package psql_test

import (
	"log"
	"os"
	"testing"

	psql "go-user-service/internal/store/postgres"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	godotenv.Load("../../../.env.dev")

	dbstr := os.Getenv("TEST_DB_URL")

	if dbstr == "" {
		log.Fatal("No database url provided")
	}

	err := psql.DropEverythingInDatabase(dbstr)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
