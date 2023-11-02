package productsUsecases

import (
	"math"

	"github.com/oiceo123/kawaii-shop-tutorial/modules/entities"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/products"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/products/productsRepositories"
)

type IProductsUsecase interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindProducts(req *products.ProductFilter) *entities.PaginateRes
	AddProduct(req *products.Product) (*products.Product, error)
}

type productsUsecase struct {
	productsRepository productsRepositories.IProductsRepository
}

func ProductsUsecase(productsRepository productsRepositories.IProductsRepository) IProductsUsecase {
	return &productsUsecase{
		productsRepository: productsRepository,
	}
}

func (u *productsUsecase) FindOneProduct(productId string) (*products.Product, error) {
	product, err := u.productsRepository.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (u *productsUsecase) FindProducts(req *products.ProductFilter) *entities.PaginateRes {
	products, count := u.productsRepository.FindProducts(req)

	return &entities.PaginateRes{
		Data:      products,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *productsUsecase) AddProduct(req *products.Product) (*products.Product, error) {
	product, err := u.productsRepository.InsertProduct(req)
	if err != nil {
		return nil, err
	}
	return product, nil
}