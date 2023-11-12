package servers

import (
	"github.com/oiceo123/kawaii-shop-tutorial/modules/products/productsHandlers"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/products/productsRepositories"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/products/productsUsecases"
)

type IProductsModule interface {
	Init()
	Repository() productsRepositories.IProductsRepository
	Usecase() productsUsecases.IProductsUsecase
	Handler() productsHandlers.IProductsHandler
}

type productsModule struct {
	*moduleFactory
	repository productsRepositories.IProductsRepository
	usecase    productsUsecases.IProductsUsecase
	handler    productsHandlers.IProductsHandler
}

func (m *moduleFactory) ProductsModule() IProductsModule {
	productsRepository := productsRepositories.ProductsRepository(m.server.db, m.server.cfg, m.FilesModule().Usecase())
	productsUsecase := productsUsecases.ProductsUsecase(productsRepository)
	productsHandler := productsHandlers.ProductsHandler(m.server.cfg, productsUsecase, m.FilesModule().Usecase())

	return &productsModule{
		moduleFactory: m,
		repository:    productsRepository,
		usecase:       productsUsecase,
		handler:       productsHandler,
	}
}

func (p *productsModule) Init() {
	router := p.router.Group("/products")

	router.Post("/", p.middlewares.JwtAuth(), p.middlewares.Authorize(2), p.handler.AddProduct)

	router.Patch("/:product_id", p.middlewares.JwtAuth(), p.middlewares.Authorize(2), p.handler.UpdateProduct)

	router.Get("/", p.middlewares.ApiKeyAuth(), p.handler.FindProducts)
	router.Get("/:product_id", p.middlewares.ApiKeyAuth(), p.handler.FindOneProduct)

	router.Delete("/:product_id", p.middlewares.JwtAuth(), p.middlewares.Authorize(2), p.handler.DeleteProduct)
}

func (p *productsModule) Repository() productsRepositories.IProductsRepository { return p.repository }
func (p *productsModule) Usecase() productsUsecases.IProductsUsecase           { return p.usecase }
func (p *productsModule) Handler() productsHandlers.IProductsHandler           { return p.handler }
