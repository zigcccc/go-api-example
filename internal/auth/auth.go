package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"go_api_example/internal/config"
)

func CreateOAuth2Config() oauth2.Config {

	cfg := config.GetConfig()

	return oauth2.Config{
		ClientID:     cfg.OAuth.ClientID,
		ClientSecret: cfg.OAuth.ClientSecret,
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"email"},
		Endpoint:     google.Endpoint,
	}
}

func GetSecretKey() string {
	cfg := config.GetConfig()

	return cfg.OAuth.SecretKey
}

var (
	OAuth2Config      = CreateOAuth2Config()
	OAuth2StateString = "randomstate"
)
