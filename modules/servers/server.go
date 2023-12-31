package servers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/oiceo123/kawaii-shop-tutorial/config"
)

type IServer interface {
	Start()
	GetServer() *server
}

type server struct {
	app *fiber.App
	cfg config.IConfig
	db  *sqlx.DB
}

func NewServer(cfg config.IConfig, db *sqlx.DB) IServer {
	return &server{
		cfg: cfg,
		db:  db,
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
	}
}

func (s *server) GetServer() *server {
	return s
}

func (s *server) Start() {
	// Middlewares
	middlewares := InitMiddlewares(s)
	s.app.Use(middlewares.Logger())
	s.app.Use(middlewares.Cors())

	// Modules
	v1 := s.app.Group("v1")
	modules := InitModule(v1, s, middlewares)

	modules.MonitorModule()
	modules.UsersModule()
	modules.AppinfoModule()
	modules.FilesModule().Init()
	modules.ProductsModule().Init()
	modules.OrdersModule()

	s.app.Use(middlewares.RouterCheck())

	// Gracful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("server is shutting down...")
		_ = s.app.Shutdown()
	}()

	// Listen to host:port
	/* fmt.Printf("server is starting on %v", s.cfg.App().Url()) */
	if err := s.app.Listen(s.cfg.App().Url()); err != nil {
		log.Panic(err)
	}
	fmt.Println("Running cleanup tasks...")
}
