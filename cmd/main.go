package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

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

	headersOk := handlers.AllowedHeaders([]string{"Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"}) //TODO change with the "production" host
	methodsOk := handlers.AllowedMethods([]string{"GET", "OPTIONS"})

	h1 := handlers.CombinedLoggingHandler(os.Stdout, router)
	h2 := handlers.CompressHandler(h1)

	log.Fatal(http.ListenAndServe(":8180", handlers.CORS(originsOk, headersOk, methodsOk)(h2)))
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

	fmt.Println(reflect.TypeOf(commerceProduct.Product.MasterData.Current.MasterVariant.AttributesText[0].Value))

	var product = structs.ProductFull{
		Id:    commerceProduct.Product.Id,
		Name:  commerceProduct.Product.MasterData.Current.Name,
		Image: commerceProduct.Product.MasterData.Current.MasterVariant.Images[0].Url,
		Price: commerceProduct.Product.MasterData.Current.MasterVariant.Price.Value.CentAmount,
		Size:  commerceProduct.Product.MasterData.Current.MasterVariant.AttributesText[0].Value.(string),
	}

	json.NewEncoder(w).Encode(product)
}
