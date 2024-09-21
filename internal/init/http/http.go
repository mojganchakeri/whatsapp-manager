package http

import (
	netHttp "net/http"

	"github.com/mojganchakeri/whatsapp-manager/internal/init/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Http interface {
	Listen() error
	ActiveLogger()
	ActiveRecover()
	SetGroup(groupName string) fiber.Router
	SetMiddlewares(group fiber.Router, middlewares []func(ctx *HttpContext) error)
	SetRoute(group fiber.Router, routes []Route)
}

type http struct {
	app     *fiber.App
	address string
}

type HttpMethod = string
type HttpContext = fiber.Ctx

const (
	Get    HttpMethod = netHttp.MethodGet
	Post   HttpMethod = netHttp.MethodPost
	Put    HttpMethod = netHttp.MethodPut
	Delete HttpMethod = netHttp.MethodDelete
)

const (
	StatusBadRequest          = netHttp.StatusBadRequest
	StatusInternalServerError = netHttp.StatusInternalServerError
	StatusOK                  = netHttp.StatusOK
)

type Route struct {
	Method  HttpMethod
	Path    string
	Handler func(*HttpContext) error
}

func New(cfg config.Config) Http {
	return &http{
		address: cfg.GetConfig().HttpServer.Address,
		app: fiber.New(fiber.Config{
			ErrorHandler: fiber.DefaultErrorHandler,
		}),
	}
}

func (h *http) Listen() error {
	return h.app.Listen(h.address)
}

func (h *http) ActiveLogger() {
	h.app.Use(logger.New())
}

func (h *http) ActiveRecover() {
	h.app.Use(recover.New())
}

func (h *http) SetGroup(groupName string) fiber.Router {
	return h.app.Group("/" + groupName)
}

func (h *http) SetMiddlewares(group fiber.Router, middlewares []func(ctx *HttpContext) error) {
	for _, middleware := range middlewares {
		group.Use(middleware)
	}
}

func (h *http) SetRoute(group fiber.Router, routes []Route) {
	for _, route := range routes {
		group.Add(string(route.Method), route.Path, route.Handler)
	}
}
