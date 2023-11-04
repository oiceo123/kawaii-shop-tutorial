package ordersUsecases

import (
	"github.com/oiceo123/kawaii-shop-tutorial/modules/orders/ordersRepositories"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/products/productsRepositories"
)

type IOrdersUsecase interface {
}

type ordersUsecase struct {
	ordersRepository   ordersRepositories.IOrdersRepository
	productsRepository productsRepositories.IProductsRepository
}

func OrdersUsecase(ordersRepository ordersRepositories.IOrdersRepository, productsRepository productsRepositories.IProductsRepository) IOrdersUsecase {
	return &ordersUsecase{
		ordersRepository:   ordersRepository,
		productsRepository: productsRepository,
	}
}
