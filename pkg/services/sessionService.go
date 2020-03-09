package services

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const Session_Token = "session_token"
const Customer_Token = "customer_token"

func GetCookieValue(r *http.Request, cookie string) string {
	c, err := r.Cookie(cookie)
	if err != nil {
		if err == http.ErrNoCookie {
			return ""
		}
	}
	return c.Value
}

func SetCookieValue(w http.ResponseWriter, r *http.Request, cookie string, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:    cookie,
		Value:   value,
		Expires: time.Now().Add(3600 * time.Second),
	})

	r.AddCookie(&http.Cookie{
		Name:    cookie,
		Value:   value,
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
	token = GetCookieValue(r, Session_Token)

	fmt.Println("Token from cookie: " + token)

	if token == "" {
		ctx := context.Background()

		conf := &clientcredentials.Config{
			ClientID:     "ve6pA6Q5mg47bR5172nls3HP",
			ClientSecret: "m3q0jGhD5p2leAD-heznIVJ5OlEkYDiI",
			Scopes:       []string{"manage_my_profile:flexy-commerce manage_customers:flexy-commerce manage_my_shopping_lists:flexy-commerce manage_my_orders:flexy-commerce create_anonymous_token:flexy-commerce manage_my_payments:flexy-commerce manage_subscriptions:flexy-commerce view_published_products:flexy-commerce"},
			TokenURL:     "https://auth.sphere.io/oauth/flexy-commerce/anonymous/token/",
		}

		authToken, _ := conf.Token(ctx)
		token = authToken.AccessToken
		SetCookieValue(w, r, Session_Token, token)
	}

	req.Header.Set("Authorization", "Bearer "+token)
}

func SetPasswordAuthToken(w http.ResponseWriter, r *http.Request, req *http.Request, username string, password string) {
	var token string
	token = GetCookieValue(r, Customer_Token)

	if token == "" {
		ctx := context.Background()

		conf := &oauth2.Config{
			ClientID:     "ve6pA6Q5mg47bR5172nls3HP",
			ClientSecret: "m3q0jGhD5p2leAD-heznIVJ5OlEkYDiI",
			Scopes:       []string{"manage_my_profile:flexy-commerce manage_my_shopping_lists:flexy-commerce manage_my_orders:flexy-commerce manage_my_payments:flexy-commerce view_published_products:flexy-commerce"},
			Endpoint:     oauth2.Endpoint{TokenURL: "https://auth.sphere.io/oauth/flexy-commerce/customers/token/"},
		}

		authToken, err := conf.PasswordCredentialsToken(ctx, username, password)

		if err != nil {
			fmt.Println("Error creating token")
			fmt.Println(err)
			return
		}
		token = authToken.AccessToken
		SetCookieValue(w, r, Customer_Token, token)
	}

	req.Header.Set("Authorization", "Bearer "+token)
}
