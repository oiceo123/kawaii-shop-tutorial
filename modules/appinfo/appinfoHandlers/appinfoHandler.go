package appinfoHandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oiceo123/kawaii-shop-tutorial/config"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/appinfo"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/appinfo/appinfoUsecases"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/entities"
	"github.com/oiceo123/kawaii-shop-tutorial/pkg/auth"
)

type appinfoHandlerErrCode string

const (
	genertaeApiKeyErr appinfoHandlerErrCode = "appinfo-001"
	findCategoryErr   appinfoHandlerErrCode = "appinfo-002"
)

type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error
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

func (h *appinfoHandler) GenerateApiKey(c *fiber.Ctx) error {
	apiKey, err := auth.NewKawaiiAuth(
		auth.ApiKey,
		h.cfg.Jwt(),
		nil,
	)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(genertaeApiKeyErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Key string `json:"key"`
		}{
			Key: apiKey.SignToken(),
		},
	).Res()
}

func (h *appinfoHandler) FindCategory(c *fiber.Ctx) error {
	req := new(appinfo.CategoryFilter)
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	category, err := h.appinfoUsecases.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, category).Res()
}
