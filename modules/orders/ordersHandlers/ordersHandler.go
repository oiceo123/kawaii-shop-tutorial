package ordersHandlers

import (
	"github.com/oiceo123/kawaii-shop-tutorial/config"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/orders/ordersUsecases"
)

type IOrdersHandler interface {
}

type ordersHandler struct {
	cfg           config.IConfig
	ordersUsecase ordersUsecases.IOrdersUsecase
}

func OrdersHandler(cfg config.IConfig, ordersUsecase ordersUsecases.IOrdersUsecase) IOrdersHandler {
	return &ordersHandler{
		cfg:           cfg,
		ordersUsecase: ordersUsecase,
	}
}
