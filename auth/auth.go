package auth

import (
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/joho/godotenv"
)

func CreateOAuth2Config() oauth2.Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return oauth2.Config{
		ClientID:     os.Getenv("OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"email"},
		Endpoint:     google.Endpoint,
	}
}

var (
	OAuth2Config      = CreateOAuth2Config()
	OAuth2StateString = "randomstate"
	SecretKey         = os.Getenv("OAUTH2_SECRET_KEY")
)
