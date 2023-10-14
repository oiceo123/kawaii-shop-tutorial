package appinfoHandlers

import (
	"github.com/oiceo123/kawaii-shop-tutorial/config"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/appinfo/appinfoUsecases"
)

type IAppinfoHandler interface {
}

type appinfoHandler struct {
	cfg             config.IConfig
	appinfoUsecases appinfoUsecases.IAppinfoUsecase
}

func AppinfoHandler(cfg config.IConfig, appinfoUsecases appinfoUsecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg:             cfg,
		appinfoUsecases: appinfoUsecases,
	}
}
