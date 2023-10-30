package products

import (
	"github.com/oiceo123/kawaii-shop-tutorial/modules/appinfo"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/entities"
)

type Product struct {
	Id          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Category    *appinfo.Category  `json:"category"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
	Price       float64            `json:"price"`
	Images      []*entities.Images `json:"images"`
}
