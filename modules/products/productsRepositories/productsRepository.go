package productsRepositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/oiceo123/kawaii-shop-tutorial/config"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/files/filesUsecases"
)

type IProductsRepository interface {
}

type productsRepository struct {
	db          *sqlx.DB
	cfg         config.IConfig
	fileUsecase filesUsecases.IFilesUsecase
}

func ProductsRepository(db *sqlx.DB, cfg config.IConfig, fileUsecase filesUsecases.IFilesUsecase) IProductsRepository {
	return &productsRepository{
		db:          db,
		cfg:         cfg,
		fileUsecase: fileUsecase,
	}
}
