package internal

import (
	"day5/internal/config"
	"fmt"
	"net/http"
)

type UserHandler struct {
	cfg config.Config
}

func NewHandler(cfg config.Config) *UserHandler {
	return &UserHandler{
		cfg: cfg,
	}
}


// so here we are able to change the port which is not good. 
// Its not abut port, suppose we need some struct to share with handler but not want to able to mutate it 
// How to solve this; 
// ------------------------------
// Answer I got: make the feild available through getter only. 
//package config
/*
type Config struct {
    Port      string
    Env       string
    DBURL     string
    APIKey    string
    JWTSecret string
}

// Read-only interface for handlers
type Reader interface {
    GetPort() string
    GetEnv() string
    GetDBURL() string
    GetAPIKey() string
    GetJWTSecret() string
}

// Implement the interface
func (c *Config) GetPort() string      { return c.Port }
func (c *Config) GetEnv() string       { return c.Env }
func (c *Config) GetDBURL() string     { return c.DBURL }
func (c *Config) GetAPIKey() string    { return c.APIKey }
func (c *Config) GetJWTSecret() string { return c.JWTSecret }

func Load() *Config {
    return &Config{
        Port:      ":8080",
        Env:       "dev",
        DBURL:     "postgres://...",
        APIKey:    "secret-key",
        JWTSecret: "jwt-secret",
    }
}
*/
// --------------------------------------- OR --------------------------------
/*
package config

type Config struct {
    // Public fields that main needs
    Port string
    Env  string
    
    // Private fields that handlers shouldn't modify
    dbURL      string
    apiKey     string
    jwtSecret  string
}

// Getters for private fields
func (c Config) DBURL() string {
    return c.dbURL
}

func (c Config) APIKey() string {
    return c.apiKey
}

func (c Config) JWTSecret() string {
    return c.jwtSecret
}

func Load() *Config {
    return &Config{
        Port:      ":8080",
        Env:       "dev",
        dbURL:     "postgres://...",
        apiKey:    "secret-key",
        jwtSecret: "jwt-secret",
    }
}
*/

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("User Created %s", h.cfg.Port)
	h.cfg.Port = ":3000"
}