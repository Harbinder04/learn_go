package config

import "os"

type Config struct {
	Port string
	Env string
}

func Load() *Config {
	port := os.Getenv("PORT")
	env := os.Getenv("ENV")

	if port == "" {
		port = ":8080"
	}

	if env == "" {
		env = "dev"
	}

	return  &Config{
		Env: env,
		Port: port,
	}
}