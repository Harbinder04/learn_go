package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseConfig dbConfig
	ServerConfig   serverConfig
	RedisConfig    RedisConfig
}

// for our scenario rdis://:pass@localhost:port
type RedisConfig struct {
	Host     string
	Password string
	Port     string
	//  redis[s]://[[username][:password]@][host][:port][/db-number]:
}

type dbConfig struct {
	host     string
	username string
	password string
	port     int64
	dbName   string
}

type serverConfig struct {
	Port              string
	Env               string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
}

func NewConfig() *Config {
	/* In go when we run 
	-go run path to main.go 
	-so if our main is not in the main and we cd into the folder then go start seeing 
	env from that point so if we are running will full path from root like here, 
	- `go run cmd/server/main.go`
	- then .env is at the same level 
	- but if we cd into `cmd/server` and then run main.go
	- env path will be "../../.env.dev"
	*/
	godotenv.Load("./.env.dev")
	env := os.Getenv("ENV")

	fmt.Print(env)
	if env == "" {
		panic(errors.New("env not provided"))
	}

	db, err := newDbConnection(env)
	if err != nil {
		panic(err)
	}

	server, err := newServerConfig()
	if err != nil {
		panic(err)
	}

	rd, err := newRedisConfig()

	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		DatabaseConfig: *db,
		ServerConfig:   *server,
		RedisConfig:    *rd,
	}
}

func newDbConnection(env string) (*dbConfig, error) {
	var prefixEnv string
	switch env {
	case "dev":
		prefixEnv = "DEV_"
	case "test":
		prefixEnv = "TEST_"
	default:
		return nil, errors.New("Unknown env provided")
	}

	dbhost := os.Getenv(prefixEnv + "DB_HOST")
	if dbhost == "" {
		return nil, errors.New("Host not provided")
	}
	databasePort, err := strconv.ParseInt(
		os.Getenv(prefixEnv+"DB_PORT"), 0, 64,
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
		host:     dbhost,
		username: databaseUsername,
		password: databasePassword,
		port:     databasePort,
		dbName:   databaseName,
	}, nil

}

func (db *dbConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable", db.host, db.username, db.password, db.port, db.dbName)
}

func newServerConfig() (*serverConfig, error) {
	env := os.Getenv("ENV")
	if env == "" {
		return &serverConfig{}, errors.New("server env not porvided")
	}
	svr_port := os.Getenv("SERVER_PORT")
	if svr_port == "" {
		return &serverConfig{}, errors.New("server port not provided")
	}

	return &serverConfig{
		Port:              svr_port,
		Env:               env,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}, nil
}

func newRedisConfig() (*RedisConfig, error) {
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		log.Fatal("redis port not provided")
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisPort == "" {
		log.Fatal("redis host not provided")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPort == "" {
		log.Fatal("redis password not provided")
	}

	return &RedisConfig{
		Host:     redisHost,
		Port:     redisPort,
		Password: redisPassword,
	}, nil
}
