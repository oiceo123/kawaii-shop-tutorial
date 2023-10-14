package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/appinfo/appinfoHandlers"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/appinfo/appinfoRepositories"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/appinfo/appinfoUsecases"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/middlewares/middlewaresHandlers"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/middlewares/middlewaresRepositories"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/middlewares/middlewaresUsecases"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/monitor/monitorHandlers"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/users/usersHandlers"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/users/usersRepositories"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/users/usersUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
}

type moduleFactory struct {
	router      fiber.Router
	server      *server
	middlewares middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, middlewares middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		router:      r,
		server:      s,
		middlewares: middlewares,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.server.cfg)

	m.router.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.server.db)
	usecase := usersUsecases.UsersUsecase(m.server.cfg, repository)
	handler := usersHandlers.UsersHandler(m.server.cfg, usecase)

	router := m.router.Group("/users")

	router.Post("/signup", m.middlewares.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.middlewares.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.middlewares.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.middlewares.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup-admin", m.middlewares.JwtAuth(), m.middlewares.Authorize(2), handler.SignUpAdmin)

	router.Get("/:user_id", m.middlewares.JwtAuth(), m.middlewares.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.middlewares.JwtAuth(), m.middlewares.Authorize(2), handler.GenerateAdminToken)

	// Initial admin ขึ้นมา 1 คน ใน Db (Insert ใน SQL)
	// Generate Admin Key
	// ทุกครั้งที่ทำการสมัคร Admin เพิ่ม ให้ส่ง Admin Token มาด้วยทุกครั้ง ผ่าน Middleware
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.server.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.server.cfg, usecase)

	router := m.router.Group("/appinfo")

	router.Get("/apikey", m.middlewares.JwtAuth(), m.middlewares.Authorize(2), handler.GenerateApiKey)
}
