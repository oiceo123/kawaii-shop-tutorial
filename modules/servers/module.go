package servers

import (
	"github.com/gofiber/fiber/v2"
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

	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", handler.RefreshPassport)
}
