package psql_test

import (
	"database/sql"
	"os"
	"testing"

	"go-user-service/config"
	psql "go-user-service/internal/store/postgres"
)

func TestMain(m *testing.M) {
	cfg := config.NewConfig()

	db, err := sql.Open("postgres", cfg.DatabaseConfig.GetConnectionString())
	err = psql.DropEverythingInDatabase(db)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
