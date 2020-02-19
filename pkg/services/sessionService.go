package services

import (
	"net/http"
	"time"
)

func GetCookieToken(r *http.Request) string {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return ""
		}
	}
	return c.Value
}

func SetCookieToken(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   token,
		Expires: time.Now().Add(3600 * time.Second),
	})

	r.AddCookie(&http.Cookie{
		Name:    "session_token",
		Value:   token,
		Expires: time.Now().Add(3600 * time.Second),
	})
}
