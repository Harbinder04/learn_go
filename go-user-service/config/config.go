package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseConfig dbConfig
}

type dbConfig struct {
	host string
	username string
	password string
	port int64
	dbName string
}

func NewConfig() *Config {
	_ = godotenv.Load("../.env.dev")
	env := os.Getenv("ENV")
	if env == "" {
		panic(errors.New("Env not provided"))
	}

	db, err := newDbConnection(env)
	if err != nil {
		panic(err)
	}

	return &Config{
		DatabaseConfig: *db,
	}
}

func newDbConnection(env string) (*dbConfig, error) {
	var prefixEnv string
	switch env {
	case "dev": 
		prefixEnv = "DEV_"
	case "test":
		prefixEnv = "TEST_"
	}

	dbhost := os.Getenv(prefixEnv + "DB_HOST")
	if dbhost == "" {
		return nil, errors.New("Host not provided")
	}
	databasePort, err := strconv.ParseInt(
		os.Getenv(prefixEnv + "DB_PORT"), 0, 64,
	)
	if err != nil {
		return nil, errors.New("could not convert db port to int")
	}
	
	databaseUsername := os.Getenv(prefixEnv + "DB_USERNAME")
	if databaseUsername == "" {
		return nil, errors.New("databaseUsername was empty")
	}
	databaseName := os.Getenv(prefixEnv + "DB_NAME")
	if databaseName == "" {
		return nil, errors.New("databaseName was empty")
	}
	databasePassword := os.Getenv(prefixEnv + "DB_PASSWORD")
	if databasePassword == "" {
		return nil, errors.New("databasePassword was empty")
	}
	return &dbConfig{
		host: dbhost,
		username:     databaseUsername,
		password: databasePassword,
		port:     databasePort,
		dbName:     databaseName,
	}, nil

}

func (db *dbConfig) getConnectionstring() string {
	return fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable", db.host, db.username, db.password, db.port, db.dbName)
}