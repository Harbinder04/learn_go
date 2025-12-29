package main

import (
	// "day5/internal"
	// "day5/internal/config"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

// dir returns the absolute path of the given environment file (envFile) in the Go module's
// root directory. It searches for the 'go.mod' file from the current working directory upwards
// and appends the envFile to the directory containing 'go.mod'.
// It panics if it fails to find the 'go.mod' file.
func dir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}

func main() {
	// err := godotenv.Load("../../.env.dev") // absolute path 
	err := godotenv.Load(dir(".env.dev"))

    if err != nil {
      log.Fatal("Error loading .env file", err)
    }
	
    r := chi.NewRouter()
    // cfg := config.Load()

	// problem: 4
	// if cfg.Port == "" || cfg.Env == "" {
	// 	log.Fatal("Port not Provided")
	// }

	// nh := internal.NewHandler(*cfg)
    env := os.Getenv("ENV")
	port := os.Getenv("PORT")

    if env == "dev" {
        r.Use(middleware.Logger)
        r.Use(middleware.RequestID)
        log.Println("Development mode: verbose logging enabled")
    }
    
	// 🪧 just my thing to see some thing Not realated to problems 
	//  read file for more info 
	//-----------------------------
	// r.Get("/port", nh.CreateUser)
	//-----------------------------

    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        if env == "prod" {
            log.Printf("Handling request: %s %s", r.Method, r.URL.Path)
        }
        w.Write([]byte("Hello World"))
    })
    
    fmt.Printf("Server started on %s (env: %s)\n", port, env)
    http.ListenAndServe(":"+port, r)
}