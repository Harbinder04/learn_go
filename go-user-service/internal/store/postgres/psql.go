package psql

import (
	"database/sql"
	"errors"
	"log"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunUpMigrations(db *sql.DB) error {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../../../migrations")
	migrationDir := filepath.Join("file://" + basePath)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	//⚠️ don't close database connection as using same connection in whole test suit.
	// defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(
		migrationDir, "postgres",
		driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No migrations to apply")
			return err
		} else {
			log.Fatal(err)
			return err
		}
	}

	//⚠️ don't close database connection as using same connection in whole test suit.
	// m.Close()

	return nil
}

func RunDownMigrations(db *sql.DB) error {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../../../migrations")
	migrationDir := filepath.Join("file://" + basePath)

	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(
		migrationDir, "postgres",
		driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No migrations to apply")
			return err
		} else {
			log.Fatal(err)
			return err
		}
	}

	m.Close()

	return nil
}

func DropEverythingInDatabase(db *sql.DB) error {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../../../migrations")
	migrationDir := filepath.Join("file://" + basePath)

	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(migrationDir, "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Drop(); err != nil {
		return err
	}
	srcErr, dbErr := m.Close()
	if srcErr != nil || dbErr != nil {
		return errors.Join(srcErr, dbErr)
	}

	return nil
}
