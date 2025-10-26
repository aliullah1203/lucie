package routes

import (
	"authentication/controllers"
	"authentication/helpers"
	"authentication/services"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func RegisterHTTPRoutes(router *mux.Router) {
	// Health Check
	router.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"pong"}`))
	}).Methods("GET")

	// Signup & Login
	router.HandleFunc("/api/signup", controllers.Signup).Methods("POST")
	router.HandleFunc("/api/login", controllers.Login).Methods("POST")
	// Logout route (protected, user must be logged in)
	router.Handle("/api/logout", helpers.AuthMiddleware(http.HandlerFunc(controllers.Logout))).Methods("POST")

	// Users routes
	router.Handle("/api/users", helpers.AuthMiddleware(http.HandlerFunc(controllers.GetUsers), "ADMIN", "SUPER_ADMIN")).Methods("GET")
	router.Handle("/api/users/{id}", helpers.AuthMiddleware(http.HandlerFunc(controllers.GetUser), "ADMIN", "SUPER_ADMIN")).Methods("GET")

	// Google OAuth login
	router.HandleFunc("/api/oauth/google/login", func(w http.ResponseWriter, r *http.Request) {
		state := generateState()
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_state",
			Value:    state,
			Path:     "/",
			HttpOnly: true,
			Secure:   isHTTPS(r),
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Now().Add(5 * time.Minute),
		})

		url := services.GetGoogleLoginURL(state)
		http.Redirect(w, r, url, http.StatusFound)
	}).Methods("GET")

	// Google OAuth callback
	router.HandleFunc("/api/oauth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		stateCookie, err := r.Cookie("oauth_state")
		if err != nil || stateCookie.Value == "" {
			http.Error(w, "missing state", http.StatusBadRequest)
			return
		}

		returnedState := r.URL.Query().Get("state")
		if returnedState != stateCookie.Value {
			http.Error(w, "invalid state", http.StatusBadRequest)
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "missing code", http.StatusBadRequest)
			return
		}

		user, token, err := services.HandleGoogleCallback(code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"` + user.ID.String() + `","email":"` + user.Email + `","token":"` + token + `"}`))
	}).Methods("GET")
}

func generateState() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	return base64.RawURLEncoding.EncodeToString(buf)
}

func isHTTPS(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto == "https" {
		return true
	}
	return false
}
