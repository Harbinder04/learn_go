package psql_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"go-user-service/config"
	internal "go-user-service/internal/store"
	psql "go-user-service/internal/store/postgres"

	"github.com/golang-migrate/migrate/v4"
)

func TestUserRepository_CreateAndGet(t *testing.T) {
	cfg := config.NewConfig()

	db, err := sql.Open("postgres", cfg.DatabaseConfig.GetConnectionString())

	if err != nil {
		t.Error(err)
	}

	storage := internal.NewSQLUserStore(db)

	err = psql.RunUpMigrations(db)
	if err != nil {
		t.Errorf("Test setup failed for: CreateUser, with err: %v", err)
	}

	t.Run("Should create a new User", func(t *testing.T) {
		newUser := internal.User{
			Id:    "B23",
			Name:  "Harbinder",
			Email: "harbinder621@gmail.com",
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, _, err = storage.Create(ctx, newUser)

		if err != nil {
			t.Errorf("Failed to create new user with err: %v", err)
			return
		}

		var queryResult internal.User

		err = db.QueryRowContext(ctx, "Select id, name, email From users Where email=$1", "harbinder621@gmail.com").Scan(&queryResult.Id, &queryResult.Name, &queryResult.Email)

		if err != nil {
			t.Errorf("This was query err: %v", err)
			return
		}

		if queryResult.Name != newUser.Name {
			t.Error(`failed 'should create a new user' wanted name did not match returned value`)
			return
		}
		if queryResult.Email != newUser.Email {
			t.Error(`failed 'should create a new user' wanted email did not match 
				returned value`)
			return
		}

		if queryResult.Id != newUser.Id {
			t.Error(`failed 'should create a new user' wanted id did not match 
				returned value`)
			return
		}
	})

	t.Cleanup(func() {
		err := psql.RunDownMigrations(db)
		if err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				return
			}
			t.Errorf("test cleanup failed for: CreateUser, with err: %v", err)
		}

	})
}
