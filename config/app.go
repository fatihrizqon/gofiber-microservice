package config

import (
	"fmt"

	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/handler"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/middleware"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/route"
	"github.com/fatihrizqon/gofiber-microservice/internal/rbac"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/service"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/fatihrizqon/gofiber-microservice/internal/util/storage"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
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

	// ── Static Files ─────────────────────────────────────────────────────────
	config.App.Get("/uploads*", static.New("./public/uploads"))

	// ── Storage ──────────────────────────────────────────────────────────────
	port := config.Config.GetInt("web.port")
	baseURL := fmt.Sprintf("http://localhost:%d/uploads", port)
	localStorage := storage.NewLocalStorage("./public/uploads", baseURL)

	// ── Repositories ──────────────────────────────────────────────────────────
	userRepository := repository.NewUserRepository(config.DB)
	authRepository := repository.NewAuthRepository(config.DB)
	tokenRepository := repository.NewTokenRepository(config.DB)
	rbacRepository := repository.NewRbacRepository(config.DB)
	fileRepository := repository.NewFileRepository(config.DB)

	// ── Email Service ─────────────────────────────────────────────────────────
	emailConfig := service.EmailConfig{
		Host:        config.Config.GetString("smtp.host"),
		Port:        config.Config.GetInt("smtp.port"),
		Username:    config.Config.GetString("smtp.username"),
		Password:    config.Config.GetString("smtp.password"),
		SenderName:  config.Config.GetString("smtp.sender_name"),
		SenderEmail: config.Config.GetString("smtp.sender_email"),
	}
	emailService := service.NewEmailService(emailConfig, config.Log)

	// ── Services ──────────────────────────────────────────────────────────────
	userService := service.NewUserService(userRepository, config.Validate)
	appURL := fmt.Sprintf("http://localhost:%d", config.Config.GetInt("web.port"))
	authService := service.NewAuthService(authRepository, tokenRepository, config.Validate, emailService, appURL)
	fileService := service.NewFileService(fileRepository, localStorage)

	// ── Handlers ──────────────────────────────────────────────────────────────
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService, config.Production)
	fileHandler := handler.NewFileHandler(fileService)

	// ── Middleware ────────────────────────────────────────────────────────────
	authMiddleware := middleware.NewAuth(tokenRepository)

	rbacEngine := rbac.New(rbac.Config{
		Store: rbacRepository,
		UserLookup: func(c fiber.Ctx) string {
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
		FileHandler:    fileHandler,
		Production:     config.Production,
	}

	routeConfig.Setup()
}
