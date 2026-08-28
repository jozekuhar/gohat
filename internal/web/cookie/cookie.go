package cookie

import (
	"net/http"
	"time"
)

const (
	cookieKeySession    = "session"
	cookieKeyOAuthState = "oauth_state"
)

func GetSession(r *http.Request) (string, error) {
	c, err := r.Cookie(cookieKeySession)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func SetSession(w http.ResponseWriter, sessionID string, expiresAt time.Time) {
	maxAge := int(time.Until(expiresAt).Seconds())
	http.SetCookie(w, &http.Cookie{
		Name:     cookieKeySession,
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   maxAge,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieKeySession,
		Value:    "",
		Path:     "/",
		Domain:   "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func GetOAuthState(r *http.Request) (string, error) {
	c, err := r.Cookie(cookieKeyOAuthState)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func SetOAuthState(w http.ResponseWriter, state string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieKeyOAuthState,
		Value:    state,
		Path:     "/",
		MaxAge:   8640, // 1 day
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
