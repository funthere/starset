package product

import "github.com/funthere/starset/orderservice/domain"

type Product struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Price       int64       `json:"price"`
	Owner       domain.User `json:"owner"`
}
