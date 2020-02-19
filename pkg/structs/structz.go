package structs

type Product struct {
	Id    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Image string `json:"image, omitempty"`
	Price int64  `json:"price, omitempty"`
}

type ProductFull struct {
	Id    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Image string `json:"image, omitempty"`
	Price int64  `json:"price, omitempty"`
	Size  string `json:"size"`
}

type Products struct {
	Products []Product `json:"products"`
}

type Cart struct {
	Id         string   `json:"id"`
	LineItems  []string `json:"lineItems"`
	TotalPrice Price    `json:"totalPrice"`
}

type Price struct {
	CurrencyCode   string `json:"currencyCode"`
	CentAmount     int64  `json:"centAmount"`
	FractionDigits int64  `json:"fractionDigits"`
}

type PriceValue struct {
	Value Price `json:"value"`
}

type AddToCartRequest struct {
	ProductId string `json:"productId"`
	Quantity  int64  `json:"quantity"`
}
