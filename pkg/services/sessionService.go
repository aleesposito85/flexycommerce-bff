package services

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

//GetCartIDAndVersion retrieve the cartId and cartVersion from the request cookie
func GetCartIDAndVersion(r *http.Request) (string, int64) {
	cartCookie, err := r.Cookie("session_cart")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", 0
		}
	}

	values := strings.Split(cartCookie.Value, ",")
	version, _ := strconv.ParseInt(values[1], 10, 64)

	return values[0], version
}

func StoreCartIDAndVersion(w http.ResponseWriter, r *http.Request, cartID string, cartVersion int64) {

	fmt.Println(cartVersion)
	fmt.Println("storing in cookie cart version: " + strconv.FormatInt(cartVersion, 10))

	value := cartID + "," + strconv.FormatInt(cartVersion, 10)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_cart",
		Value:   value,
		Expires: time.Now().Add(3600 * time.Second),
	})

	r.AddCookie(&http.Cookie{
		Name:    "session_cart",
		Value:   value,
		Expires: time.Now().Add(3600 * time.Second),
	})
}
