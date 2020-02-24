package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"../pkg/commercetools"
	"../pkg/structs"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//ChatbotRequest is the request incoming from the website chat
type ChatbotRequest struct {
	Text      string `json:"text"`
	SessionID string `json:"sessionId"`
	ProjectID string `json:"projectId"`
}

// main function to boot up everything
func main() {
	log.Println("Starting the listeners on port 8180 for BFF")
	router := mux.NewRouter()

	router.HandleFunc("/getProducts", getProducts).Methods("GET")
	router.HandleFunc("/product", getProduct).Methods("GET")
	router.HandleFunc("/cart", getCart).Methods("GET")
	router.HandleFunc("/addToCart", addToCart).Methods("POST")
	router.HandleFunc("/updateCartItem", updateCartItem).Methods("POST")

	headersOk := handlers.AllowedHeaders([]string{"Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"}) //TODO change with the "production" host
	methodsOk := handlers.AllowedMethods([]string{"GET", "OPTIONS", "POST"})
	allowCredentials := handlers.AllowCredentials()

	h1 := handlers.CombinedLoggingHandler(os.Stdout, router)
	h2 := handlers.CompressHandler(h1)

	log.Fatal(http.ListenAndServe(":8180", handlers.CORS(originsOk, headersOk, methodsOk, allowCredentials)(h2)))
}

func getProducts(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	commerceProducts := commercetools.GetProducts()

	var productSlice = []structs.Product{}

	for _, s := range commerceProducts.Products.Results {
		var product = structs.Product{
			Id:    s.Id,
			Name:  s.MasterData.Current.Name,
			Image: s.MasterData.Current.MasterVariant.Images[0].Url,
			Price: s.MasterData.Current.MasterVariant.Price.Value.CentAmount,
		}
		productSlice = append(productSlice, product)
	}

	var products = structs.Products{
		Products: productSlice,
	}

	json.NewEncoder(w).Encode(products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("productId")

	w.Header().Set("Content-Type", "application/json")

	commerceProduct := commercetools.GetProduct(productID, []string{"size", "season"}, []string{"gender", "madeInItaly"})

	//fmt.Println(reflect.TypeOf(commerceProduct.Product.MasterData.Current.MasterVariant.AttributesText[0].Value))

	var product = structs.ProductFull{
		Id:    commerceProduct.Product.Id,
		Name:  commerceProduct.Product.MasterData.Current.Name,
		Image: commerceProduct.Product.MasterData.Current.MasterVariant.Images[0].Url,
		Price: commerceProduct.Product.MasterData.Current.MasterVariant.Price.Value.CentAmount,
		Size:  commerceProduct.Product.MasterData.Current.MasterVariant.AttributesText[0].Value.(string),
	}

	json.NewEncoder(w).Encode(product)
}

func getCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	commerceCart := commercetools.GetCart(w, r)

	if len(commerceCart.ID) > 0 {
		json.NewEncoder(w).Encode(commerceCart)
	} else {
		json.NewEncoder(w).Encode(nil)
	}
}

func addToCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var cartRequest structs.AddToCartRequest

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&cartRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newCart := commercetools.AddToCart(w, r, cartRequest)

	json.NewEncoder(w).Encode(newCart)
}

func updateCartItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var cartItemRequest structs.UpdateCartItemRequest

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&cartItemRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newCart := commercetools.UpdateItemToCart(w, r, cartItemRequest)

	json.NewEncoder(w).Encode(newCart)
}
