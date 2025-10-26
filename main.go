package main

import (
	"api/configs"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// InitGoogleOAuthFromEnv initializes Google OAuth configuration from environment variables.
// This is a local shim to avoid depending on an external "api/internal/services" package.
// Replace this with the real implementation or re-add the correct module import when available.
func InitGoogleOAuthFromEnv() error {
	// No-op for now; validate environment variables here if needed.
	return nil
}

func main() {
	// Load .env file (no fatal if not present, but recommended)
	// Initialize Google OAuth config from env
	if err := InitGoogleOAuthFromEnv(); err != nil {
		log.Fatalf("Google OAuth config error: %v", err)
	}
	// Initialize database
	configs.ConnectPostgres()

	// Initialize Google OAuth config from env
	if err := services.InitGoogleOAuthFromEnv(); err != nil {
		log.Fatalf("Google OAuth config error: %v", err)
	}

	// Create net/http mux and register routes
	router := mux.NewRouter()
	// Register application routes here.
	// If you have a package that provides RegisterHTTPRoutes, import it (for example "api/router")
	// and call routerPkg.RegisterHTTPRoutes(router).
	// Example:
	//    routerPkg.RegisterHTTPRoutes(router)
	// For now provide a simple root handler to avoid 404s:
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Port fallback
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server run on: %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
