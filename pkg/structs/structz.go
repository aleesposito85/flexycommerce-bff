package structs

type Product struct {
	Id    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Image string `json:"image, omitempty"`
	Price int64  `json:"price, omitempty"`
}

type Products struct {
	Products []Product `json:"products"`
}
