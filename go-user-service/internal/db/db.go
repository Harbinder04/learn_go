package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func NewdbConnection(url string) *sql.DB {
	dbCon, err := sql.Open("postgres", url);
	if err != nil {
		fmt.Print("Db connection failed")
	}

    dbCon.SetMaxOpenConns(25)
    dbCon.SetMaxIdleConns(5)
    dbCon.SetConnMaxLifetime(5 * time.Minute)
    dbCon.SetConnMaxIdleTime(5 * time.Minute)
	
	err = dbCon.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
	return dbCon
}