package productsUsecases

import (
	"github.com/oiceo123/kawaii-shop-tutorial/modules/products"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/products/productsRepositories"
)

type IProductsUsecase interface {
	FindOneProduct(productId string) (*products.Product, error)
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
