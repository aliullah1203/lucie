package api

import (
	"authentication/controllers"
	"authentication/helpers"
	"authentication/services"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

func RegisterHTTPRoutes(router *http.ServeMux) {
	// Health Check
	router.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"pong"}`))
	})

	// Signup & Login
	router.HandleFunc("/api/signup", controllers.Signup)
	router.HandleFunc("/api/login", controllers.Login)
	// Logout route (protected, user must be logged in)
	router.Handle("/api/logout", helpers.AuthMiddleware(http.HandlerFunc(controllers.Logout)))

	// Users routes
	router.Handle("/api/users", helpers.AuthMiddleware(http.HandlerFunc(controllers.GetUsers), "ADMIN", "SUPER_ADMIN"))
	router.Handle("/api/users/{id}", helpers.AuthMiddleware(http.HandlerFunc(controllers.GetUser), "ADMIN", "SUPER_ADMIN"))

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
	})

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
	})
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
