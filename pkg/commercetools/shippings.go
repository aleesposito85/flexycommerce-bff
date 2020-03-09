package commercetools

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"../services"
	"../structs"

	"fmt"
)

type CommerceShippingsResponse struct {
	structs.CommerceResultsResponse
	Results []CommerceShippingResponse `json:"results"`
}

type CommerceShippingResponse struct {
	ID          string             `json:"id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	ZoneRates   []CommerceZoneRate `json:"zoneRates"`
}

type CommerceZoneRate struct {
	Zone          interface{}            `json:"zone"`
	ShippingRates []CommerceShippingRate `json:"shippingRates"`
}

type CommerceShippingRate struct {
	Price structs.Price `json:"price"`
}

//GetShippingMethods returns the available shipping methods in the system
func GetShippingMethods(w http.ResponseWriter, r *http.Request) []CommerceShippingResponse {

	var shippingMethods []CommerceShippingResponse

	//httpClient := getHttpClient()
	url := "https://api.sphere.io/flexy-commerce/shipping-methods?country=US"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	services.SetAuthToken(w, r, req)

	response, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	if response.StatusCode > 299 {
		b, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(b))
	} else {
		json.NewDecoder(response.Body).Decode(&shippingMethods)
	}

	return shippingMethods
}
