package services

import (
	"authentication/config"
	"authentication/helpers"
	models "authentication/user"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config

func InitGoogleOAuthFromEnv() error {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/api/oauth/google/callback"
	}

	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("missing GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET")
	}

	googleOauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	return nil
}

func GetGoogleLoginURL(state string) string {
	return googleOauthConfig.AuthCodeURL(state)
}

func HandleGoogleCallback(code string) (*models.User, string, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange code for token: %v", err)
	}

	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user info: %v", err)
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, "", fmt.Errorf("failed to decode user info: %v", err)
	}

	var user models.User
	err = config.DB.Get(&user, "SELECT * FROM users WHERE email=$1", googleUser.Email)
	if err != nil {
		// Create new user
		user = models.User{
			ID:                 uuid.New(),
			Name:               googleUser.Name,
			Email:              googleUser.Email,
			Role:               "CUSTOMER",
			Status:             "ACTIVE",
			SubscriptionStatus: "SUBSCRIBED",
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		_, err := config.DB.NamedExec(`INSERT INTO users (id, name, email, role, status, subscription_status, created_at, updated_at)
			VALUES (:id, :name, :email, :role, :status, :subscription_status, :created_at, :updated_at)`, &user)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create user: %v", err)
		}
	}

	accessToken, err := helpers.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate JWT: %v", err)
	}

	return &user, accessToken, nil
}
