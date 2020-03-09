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

type CommerceSigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CommerceSigninResponse struct {
	Customer CommerceCustomerResponse `json:"customer,omitempty"`
	Cart     interface{}              `json:"cart,omitempty"`
}

type CommerceCustomerResponse struct {
	ID        string `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

func CustomerSignin(w http.ResponseWriter, r *http.Request, request structs.SigninRequest) CommerceSigninResponse {

	var signinResponse CommerceSigninResponse

	//httpClient := getHttpClient()
	url := "https://api.sphere.io/flexy-commerce/login"

	commerceToolsRequest := CommerceSigninRequest{
		Email:    request.Username,
		Password: request.Password,
	}

	fmt.Println(commerceToolsRequest.Email)
	fmt.Println(commerceToolsRequest.Password)

	data, _ := json.Marshal(commerceToolsRequest)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	services.SetPasswordAuthToken(w, r, req, request.Username, request.Password)

	response, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	if response.StatusCode > 199 {
		b, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(b))
	} else {
		json.NewDecoder(response.Body).Decode(&signinResponse)
	}

	return signinResponse
}

func GetCustomerFromPasswordFlow(w http.ResponseWriter, r *http.Request, request structs.SigninRequest) CommerceCustomerResponse {

	var customerResponse CommerceCustomerResponse

	//httpClient := getHttpClient()
	url := "https://api.sphere.io/flexy-commerce/me"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	services.SetPasswordAuthToken(w, r, req, request.Username, request.Password)

	response, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	if response.StatusCode > 299 {
		b, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(b))
	} else {
		json.NewDecoder(response.Body).Decode(&customerResponse)
	}

	return customerResponse
}
