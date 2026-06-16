package config

import (
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/handler"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/middleware"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/route"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/service"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/fatihrizqon/gofiber-rbac"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	Cors       *CORSConfig
	DB         *gorm.DB
	App        *fiber.App
	Log        *logrus.Logger
	Validate   *validator.Validate
	Config     *viper.Viper
	JWT        *JWTService
	Production bool
}

func Bootstrap(config *BootstrapConfig) {
	config.App.Use(config.Cors.Handler())

	// ── Repositories ──────────────────────────────────────────────────────────
	userRepository := repository.NewUserRepository(config.DB)
	authRepository := repository.NewAuthRepository(config.DB)
	tokenRepository := repository.NewTokenRepository(config.DB)
	rbacRepository := repository.NewRbacRepository(config.DB)

	// ── Services ──────────────────────────────────────────────────────────────
	userService := service.NewUserService(userRepository, config.Validate)
	authService := service.NewAuthService(authRepository, tokenRepository, config.Validate)

	// ── Handlers ──────────────────────────────────────────────────────────────
	userHandler := handler.NewUserHandler(userService)

	authHandler := handler.NewAuthHandler(authService, config.Production)

	// ── Middleware ────────────────────────────────────────────────────────────
	authMiddleware := middleware.NewAuth(tokenRepository)

	rbacEngine := rbac.New(rbac.Config{
		Store: rbacRepository,
		UserLookup: func(c *fiber.Ctx) string {
			claims, ok := c.Locals("auth").(*util.Claims)
			if !ok {
				return ""
			}
			return claims.UserID.String()
		},
	})

	// ── Routes ────────────────────────────────────────────────────────────────
	routeConfig := route.RouteConfig{
		App:            config.App,
		AuthMiddleware: authMiddleware,
		RbacEngine:     rbacEngine,
		UserHandler:    userHandler,
		AuthHandler:    authHandler,
		Production:     config.Production,
	}

	routeConfig.Setup()
}
