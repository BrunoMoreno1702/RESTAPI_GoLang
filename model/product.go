package model

type Product struct {
	Idq int `json:"id_product"`
	Name string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductRequest struct {
	Name string  `json:"name"`
	Price float64 `json:"price"`
}