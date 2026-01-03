package db

import (
	"database/sql"
	"log/slog"
	"time"

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
