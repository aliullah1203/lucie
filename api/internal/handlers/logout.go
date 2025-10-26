package handlers

import (
	"encoding/json"
	"net/http"
)

// Logout simply informs the client to delete the token
func Logout(w http.ResponseWriter, r *http.Request) {
	// Best practice is to clear an access token cookie if you use one
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Delete the cookie immediately
		HttpOnly: true,
	})

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}
