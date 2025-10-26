package main

import (
	"api/configs"
	"api/internal/services"
	"api/pkg/router"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (no fatal if not present, but recommended)
	// This loads the DB credentials, JWT secret, and Google OAuth keys.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, rely on env variables")
	}

	// Initialize database connection (uses variables from .env)
	configs.ConnectPostgres()

	// Initialize Google OAuth config (uses variables from .env)
	if err := services.InitGoogleOAuthFromEnv(); err != nil {
		log.Fatalf("Google OAuth config error: %v", err)
	}

	// Create net/http mux and register routes
	// Note: The package containing RegisterHTTPRoutes is 'router' (api/pkg/router).
	r := mux.NewRouter()
	router.RegisterHTTPRoutes(r)

	// Port fallback (uses PORT from .env, or defaults to 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port: %s", port)
	// Start the HTTP server
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
