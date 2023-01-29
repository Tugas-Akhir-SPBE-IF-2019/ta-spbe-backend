package auth

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"ta-spbe-backend/config"
	"golang.org/x/oauth2"
)

func Google(cfgOAuth config.OAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oAuthConfig := newOAuthConfig(cfgOAuth, SCOPE_EMAIL, SCOPE_PROFILE)

		state := uuid.NewString()

		url := oAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

		cookie := &http.Cookie{
			Name:     "state",
			Value:    state,
			Expires:  time.Now().Add(time.Minute),
			HttpOnly: true,
		}

		http.SetCookie(w, cookie)

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
