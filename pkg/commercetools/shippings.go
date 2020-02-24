package commercetools

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"../structs"

	"fmt"
)

type ShippingsResponse struct {
	structs.CommerceResultsResponse
	Results []ShippingResponse `json:"results"`
}

type ShippingResponse struct {
	ID          string             `json:"id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	Price       structs.PriceValue `json:"price"`
}

//GetCart returns the active cart of the user based on access token
func GetShippingMethods(w http.ResponseWriter, r *http.Request) CartResponse {

	var carts CartsResponse

	//httpClient := getHttpClient()
	url := "https://api.sphere.io/flexy-commerce/shipping-methods"

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
	json.NewDecoder(response.Body).Decode(&carts)

	if len(carts.Results) > 0 {
		return carts.Results[0]
	}

	return CartResponse{}
}
