package services

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2/clientcredentials"
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

func SetAuthToken(w http.ResponseWriter, r *http.Request, req *http.Request) {
	var token string
	token = GetCookieToken(r)

	fmt.Println("Token from cookie: " + token)

	if token == "" {
		ctx := context.Background()
		conf := &clientcredentials.Config{
			ClientID:     "EN4wn2L0gEUyEhim76SHs4N0",
			ClientSecret: "nCepZvUBjqL0Cw36U1t6QeWZmyzzzaLr",
			Scopes:       []string{"manage_my_profile:flexy-commerce", "manage_my_shopping_lists:flexy-commerce", "manage_my_orders:flexy-commerce", "create_anonymous_token:flexy-commerce", "manage_my_payments:flexy-commerce", "view_published_products:flexy-commerce"},
			TokenURL:     "https://auth.sphere.io/oauth/flexy-commerce/anonymous/token/",
		}

		authToken, _ := conf.Token(ctx)
		token = authToken.AccessToken
		SetCookieToken(w, r, token)

		fmt.Println("Token from new: " + token)
	}

	req.Header.Set("Authorization", "Bearer "+token)
}
