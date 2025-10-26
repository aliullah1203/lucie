package controllers

import (
	"authentication/configs"
	"authentication/internal/helpers"
	"authentication/internal/services"
	"authentication/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ====================== SIGNUP ======================
func Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	// Hash password
	hashed, err := helpers.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = hashed
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Role = "USER"
	user.Status = "ACTIVE"
	user.SubscriptionStatus = "INACTIVE"

	// Insert user
	query := `INSERT INTO users (id, name, email, phone, address, role, status, subscription_status, password, created_at, updated_at)
	          VALUES (:id,:name,:email,:phone,:address,:role,:status,:subscription_status,:password,:created_at,:updated_at)`
	_, err = configs.DB.NamedExec(query, &user)
	if err != nil {
		http.Error(w, "email already exists or db error", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(Response{Message: "Signup successful", Data: user})
}

// ====================== LOGIN ======================
func Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	var user models.User
	err := configs.DB.Get(&user, "SELECT * FROM users WHERE email=$1 AND deleted_at IS NULL", req.Email)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if !helpers.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := helpers.GenerateJWT(user.ID.String(), user.Role)
	if err != nil {
		http.Error(w, "error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(Response{Message: "Login successful", Data: map[string]string{"token": token}})
}

// ====================== LOGOUT ======================
// For stateless JWT, logout is handled client-side by deleting token.
// Optionally you can implement token blacklist in DB or Redis.
func Logout(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Response{Message: "Logout successful"})
}

// ====================== PROFILE ======================
func Profile(w http.ResponseWriter, r *http.Request) {
	// User info extracted from JWT in middleware (optional enhancement)
	w.Write([]byte("Protected profile route"))
}

// ====================== GOOGLE OAUTH CALLBACK ======================
func GoogleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "code not found", http.StatusBadRequest)
		return
	}

	userInfo, err := services.GetGoogleUserInfo(code)
	if err != nil {
		http.Error(w, "failed to get user info", http.StatusInternalServerError)
		return
	}

	var user models.User
	err = configs.DB.Get(&user, "SELECT * FROM users WHERE email=$1", userInfo.Email)
	if err != nil {
		// If user doesn't exist, create new
		user.ID = uuid.New()
		user.Name = userInfo.Name
		user.Email = userInfo.Email
		user.Role = "USER"
		user.Status = "ACTIVE"
		user.SubscriptionStatus = "INACTIVE"
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		query := `INSERT INTO users (id, name, email, role, status, subscription_status, created_at, updated_at)
		          VALUES (:id,:name,:email,:role,:status,:subscription_status,:created_at,:updated_at)`
		_, err = configs.DB.NamedExec(query, &user)
		if err != nil {
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}
	}

	token, err := helpers.GenerateJWT(user.ID.String(), user.Role)
	if err != nil {
		http.Error(w, "error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(Response{Message: "Login via Google successful", Data: map[string]string{"token": token}})
}
