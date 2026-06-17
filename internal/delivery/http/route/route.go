package route

import (
	_ "github.com/fatihrizqon/gofiber-microservice/docs"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/handler"
	"github.com/fatihrizqon/gofiber-microservice/internal/rbac"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type RouteConfig struct {
	App               *fiber.App
	AuthMiddleware    fiber.Handler
	CompanyMiddleware fiber.Handler
	RbacEngine        *rbac.RBAC
	UserHandler       *handler.UserHandler
	AuthHandler       *handler.AuthHandler
	FileHandler       *handler.FileHandler
	Production        bool
}

func (rc *RouteConfig) Setup() {
	rc.SetupGuestRoute()
	rc.SetupAuthRoute()
}

func (rc *RouteConfig) SetupGuestRoute() {
	rc.App.Post("/api/v1/auth/login", rc.AuthHandler.Login)
	rc.App.Post("/api/v1/auth/refresh", rc.AuthHandler.Refresh)
	rc.App.Post("/api/v1/users", rc.UserHandler.Create)
	rc.App.Get("/swagger/*", swagger.HandlerDefault)
}

func (rc *RouteConfig) SetupAuthRoute() {
	rc.App.Use(rc.AuthMiddleware)

	rc.App.Post("/api/v1/auth/logout", rc.AuthHandler.Logout)

	rc.App.Get("/api/v1/users", rc.RbacEngine.Require("users.read"), rc.UserHandler.FindAll)
	rc.App.Get("/api/v1/users/:id", rc.RbacEngine.Require("users.read"), rc.UserHandler.FindById)
	rc.App.Put("/api/v1/users/:id", rc.RbacEngine.Require("users.write"), rc.UserHandler.Update)
	rc.App.Delete("/api/v1/users/:id", rc.RbacEngine.Require("users.write"), rc.UserHandler.Delete)

	rc.App.Post("/api/v1/files/upload", rc.FileHandler.Upload)
}
