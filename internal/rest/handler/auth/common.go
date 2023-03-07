package auth

import (
	"time"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type TokenItem struct {
	Token     string    `json:"token"`
	Scheme    string    `json:"scheme"`
	ExpiresAt time.Time `json:"expires_at"`
}

const (
	SCOPE_EMAIL   = "https://www.googleapis.com/auth/userinfo.email"
	SCOPE_PROFILE = "https://www.googleapis.com/auth/userinfo.profile"
)

func newOAuthConfig(cfg config.OAuth, scopes ...string) oauth2.Config {
	return oauth2.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       scopes,
	}
}
