package session

import (
	"net/http"
	"time"
)

const sessionCookieKey = "session"

func GetCookie(r *http.Request) (string, error) {
	c, err := r.Cookie(sessionCookieKey)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func SetCookie(w http.ResponseWriter, sessionID string, expiresAt time.Time) {
	maxAge := int(time.Until(expiresAt).Seconds())
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieKey,
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   maxAge,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieKey,
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
