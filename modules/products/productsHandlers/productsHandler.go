package productsHandlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/oiceo123/kawaii-shop-tutorial/config"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/entities"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/files/filesUsecases"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/products/productsUsecases"
)

type productsHandlerErrCode string

const (
	findOneProductErr productsHandlerErrCode = "product-001"
)

type IProductsHandler interface {
	FindOneProduct(c *fiber.Ctx) error
}

type productsHandler struct {
	cfg             config.IConfig
	productsUsecase productsUsecases.IProductsUsecase
	filesUsecase    filesUsecases.IFilesUsecase
}

func ProductsHandler(cfg config.IConfig, productsUsecase productsUsecases.IProductsUsecase, filesUsecase filesUsecases.IFilesUsecase) IProductsHandler {
	return &productsHandler{
		cfg:             cfg,
		productsUsecase: productsUsecase,
		filesUsecase:    filesUsecase,
	}
}

func (h *productsHandler) FindOneProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productsUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneProductErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}
