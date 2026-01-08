package db

import (
	"database/sql"
	"log"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewdbConnection(logger *slog.Logger, url string) (sq *sql.DB, err error) {
	dbCon, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	dbCon.SetMaxOpenConns(25)
	dbCon.SetMaxIdleConns(5)
	dbCon.SetConnMaxLifetime(5 * time.Minute)
	dbCon.SetConnMaxIdleTime(5 * time.Minute)

	// This won't work because it is being ignored or overridden by the os's TCP connection timeout as error is:
	// ❌ dial tcp 127.0.0.1:5432: i/o timeout
	// Solution is add the timeout in db url

	// we need this to prevent the infinite ping request
	// this will cancel after 10s and if the connection is successfull it will get cancelled.

	//------------------------------------------------
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	//------------------------------------------------

	err = dbCon.Ping()
	if err != nil {
		return nil, err
	}

	logger.Info("Successfully connected to PostgreSQL!")
	return dbCon, nil
}

//⚠️todo: remove this later
func ApplyMigrations(db string) (error) {
	m, err := migrate.New(
		"file://../../migrations",
		db)
	if err != nil {
		log.Fatal(err)
	}
	// apply .down migrations according to the version
	// if err :=  m.Down(); err != nil {
	//     if err == migrate.ErrNoChange {
	//         log.Println("No migrations to apply")
	//     } else {
	//         log.Fatal(err)
	//     }
	// }

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No migrations to apply")
			return err
		} else {
			log.Fatal(err)
			return err
		}
	}
	// scheme version
	scheme_version, _, _ := m.Version()
	log.Printf("Scheme version is :%d", scheme_version)

	return nil
}
