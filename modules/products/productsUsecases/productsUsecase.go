package productsUsecases

import "github.com/oiceo123/kawaii-shop-tutorial/modules/products/productsRepositories"

type IProductsUsecase interface {
}

type productsUsecase struct {
	productsRepository productsRepositories.IProductsRepository
}

func ProductsUsecase(productsRepository productsRepositories.IProductsRepository) IProductsUsecase {
	return &productsUsecase{
		productsRepository: productsRepository,
	}
}
