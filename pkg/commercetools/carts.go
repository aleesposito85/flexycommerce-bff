package commercetools

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"../services"
	"../structs"

	"fmt"

	"golang.org/x/oauth2/clientcredentials"
)

type AddToCartRequest struct {
	Version int64              `json:"version"`
	Actions []CartUpdateAction `json:"actions"`
}

type CartUpdateAction struct {
	Action     string `json:"action,omitempty"`
	ProductId  string `json:"productId,omitempty"`
	LineItemId string `json:"lineItemId,omitempty"`
	Quantity   int64  `json:"quantity"`
}

type CartsResponse struct {
	Limit   int64          `json:"limit"`
	Offset  int64          `json:"offset"`
	Count   int64          `json:"count"`
	Total   int64          `json:"total"`
	Results []CartResponse `json:"results"`
}

type CartResponse struct {
	ID          string        `json:"id,omitempty"`
	Version     int64         `json:"version,omitempty"`
	AnonymousId string        `json:"anonymousId,omitempty"`
	TotalPrice  structs.Price `json:"totalPrice,omitempty"`
	LineItems   []LineItem    `json:"lineItems,omitempty"`
}

type LineItem struct {
	Id         string             `json:"id"`
	Name       interface{}        `json:"name"`
	ProductId  string             `json:"productId"`
	Price      structs.PriceValue `json:"price"`
	Quantity   int64              `json:"quantity"`
	TotalPrice structs.Price      `json:"totalPrice"`
	Variant    Variant            `json:"variant"`
}

type Variant struct {
	Images []Image `json:"images"`
}

type Image struct {
	Url string `json:"url"`
}

type CreateCartRequest struct {
	Currency string `json:"currency"`
}

var httpClient = getHTTPClient()

//GetCart returns the active cart of the user based on access token
func GetCart(w http.ResponseWriter, r *http.Request) CartResponse {

	var carts CartsResponse

	//httpClient := getHttpClient()
	url := "https://api.sphere.io/flexy-commerce/me/carts"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	setAuthToken(w, r, req)

	response, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	fmt.Println("Get cart response status: " + response.Status)

	if response.StatusCode > 299 {
		b, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(b))
	}

	json.NewDecoder(response.Body).Decode(&carts)

	if len(carts.Results) > 0 {
		return carts.Results[0]
	}

	return CartResponse{}
}

//AddToCart adds a product with the specified quantity to the user cart (new or existing)
func AddToCart(w http.ResponseWriter, r *http.Request, request structs.AddToCartRequest) CartResponse {

	cartID, cartVersion := services.GetCartIDAndVersion(r)

	if cartID == "" {
		cart := GetOrCreateCart(w, r)
		cartID = cart.ID
		cartVersion = cart.Version
	}

	url := "https://api.sphere.io/flexy-commerce/me/carts/" + cartID

	fmt.Println("Calling url: " + url)

	action := CartUpdateAction{
		Action:    "addLineItem",
		ProductId: request.ProductId,
		Quantity:  request.Quantity,
	}

	commerceToolsRequest := AddToCartRequest{
		Version: cartVersion,
		Actions: []CartUpdateAction{action},
	}
	data, _ := json.Marshal(commerceToolsRequest)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	setAuthToken(w, r, req)

	// Send request
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()
	fmt.Println("Add to cart response status: " + resp.Status)

	if resp.StatusCode > 299 {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
	}

	var cart CartResponse
	json.NewDecoder(resp.Body).Decode(&cart)

	services.StoreCartIDAndVersion(w, r, cart.ID, cart.Version)
	return cart
}

//UpdateItemToCart updates the quantity of an item in the cart
func UpdateItemToCart(w http.ResponseWriter, r *http.Request, request structs.UpdateCartItemRequest) CartResponse {
	cartID, cartVersion := services.GetCartIDAndVersion(r)

	if cartID == "" {
		cart := GetOrCreateCart(w, r)
		cartID = cart.ID
		cartVersion = cart.Version
	}

	url := "https://api.sphere.io/flexy-commerce/me/carts/" + cartID

	fmt.Println("Calling url: " + url)

	action := CartUpdateAction{
		Action:     "changeLineItemQuantity",
		LineItemId: request.ItemId,
		Quantity:   request.Quantity,
	}

	commerceToolsRequest := AddToCartRequest{
		Version: cartVersion,
		Actions: []CartUpdateAction{action},
	}
	data, _ := json.Marshal(commerceToolsRequest)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	setAuthToken(w, r, req)

	// Send request
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()
	fmt.Println("Add to cart response status: " + resp.Status)

	if resp.StatusCode > 299 {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
	}
	var updatedCart CartResponse
	json.NewDecoder(resp.Body).Decode(&updatedCart)

	services.StoreCartIDAndVersion(w, r, updatedCart.ID, updatedCart.Version)

	return updatedCart
}

//GetOrCreateCart used to fetch an active cart based on user access token or create a new one if none is available
func GetOrCreateCart(w http.ResponseWriter, r *http.Request) CartResponse {
	cart := GetCart(w, r)

	if cart.Version == 0 {
		cart = CreateCart(w, r)
	}

	services.StoreCartIDAndVersion(w, r, cart.ID, cart.Version)

	return cart
}

//CreateCart creates a new cart regardless of the already existing ones
func CreateCart(w http.ResponseWriter, r *http.Request) CartResponse {
	//httpClient := getHttpClient()
	url := "https://api.sphere.io/flexy-commerce/me/carts"

	request := CreateCartRequest{
		Currency: "USD",
	}

	data, _ := json.Marshal(request)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	setAuthToken(w, r, req)

	// Send request
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()
	fmt.Println("Create cart response status: " + resp.Status)
	var cart CartResponse
	json.NewDecoder(resp.Body).Decode(&cart)
	return cart
}

func getHTTPClient() *http.Client {
	return &http.Client{}
}

func setAuthToken(w http.ResponseWriter, r *http.Request, req *http.Request) {
	var token string
	token = services.GetCookieToken(r)

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
		services.SetCookieToken(w, r, token)

		fmt.Println("Token from new: " + token)
	}

	req.Header.Set("Authorization", "Bearer "+token)
}
