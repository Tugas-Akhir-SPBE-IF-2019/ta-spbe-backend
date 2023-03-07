package auth

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func (handler *authHandler) Google(w http.ResponseWriter, r *http.Request) {
	oAuthConfig := newOAuthConfig(handler.cfgOAuth, SCOPE_EMAIL, SCOPE_PROFILE)

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
