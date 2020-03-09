package structs

type Product struct {
	Id    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Image string `json:"image, omitempty"`
	Price int64  `json:"price, omitempty"`
}

type ProductFull struct {
	Id          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image, omitempty"`
	Price       int64  `json:"price, omitempty"`
	Size        string `json:"size"`
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

type UpdateCartItemRequest struct {
	ItemId   string `json:"itemId"`
	Quantity int64  `json:"quantity"`
}

type Address struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	StreetName   string `json:"streetName"`
	StreetNumber string `json:"streetNumber"`
	PostalCode   string `json:"postalCode"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
}

type CommerceResultsResponse struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
	Count  int64 `json:"count"`
	Total  int64 `json:"total"`
}

type ShippingMethod struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Price       Price  `json:"price"`
}

type ShippingsResponse struct {
	Shippings []ShippingMethod `json:"shippings"`
}

type SigninRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
