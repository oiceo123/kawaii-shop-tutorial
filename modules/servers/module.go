package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oiceo123/kawaii-shop-tutorial/modules/monitor/monitorHandlers"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	router fiber.Router
	server *server
}

func InitModule(r fiber.Router, s *server) IModuleFactory {
	return &moduleFactory{
		router: r,
		server: s,
	}
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.server.cfg)

	m.router.Get("/", handler.HealthCheck)
}
