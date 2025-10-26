package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"api/configs"
	"api/internal/helpers"
	"api/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.JSON(w, map[string]string{"error": "invalid JSON"}, http.StatusBadRequest)
		return
	}

	// 1. Check duplicate email or phone
	var count int
	err := configs.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE email=$1 OR phone=$2", user.Email, user.Phone)
	if err != nil {
		helpers.JSON(w, map[string]string{"error": "database error"}, http.StatusInternalServerError)
		return
	}
	if count > 0 {
		helpers.JSON(w, map[string]string{"error": "email or phone already exists"}, http.StatusBadRequest)
		return
	}

	// 2. Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.JSON(w, map[string]string{"error": "failed to hash password"}, http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// 3. Fill other fields
	user.ID = uuid.New()
	user.Role = "CUSTOMER"
	user.Status = "ACTIVE"
	user.SubscriptionStatus = "SUBSCRIBED"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// 4. Insert user into DB
	_, err = configs.DB.NamedExec(`INSERT INTO users 
        (id, name, email, phone, address, role, status, subscription_status, password, created_at, updated_at) 
        VALUES 
        (:id,:name,:email,:phone,:address,:role,:status,:subscription_status,:password,:created_at,:updated_at)`, &user)
	if err != nil {
		helpers.JSON(w, map[string]string{"error": "insert error"}, http.StatusInternalServerError)
		return
	}

	// 5. Generate JWT token
	token, err := helpers.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		helpers.JSON(w, map[string]string{"error": "token generation failed"}, http.StatusInternalServerError)
		return
	}

	// 6. Send JSON response
	helpers.JSON(w, map[string]interface{}{
		"message": "User created successfully",
		"token":   token,
	}, http.StatusCreated)
}
