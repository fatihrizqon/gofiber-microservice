package config

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
)

type CORSConfig struct {
	AllowOrigins     string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
}

func (c *CORSConfig) Handler() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     c.AllowOrigins,
		AllowCredentials: c.AllowCredentials,
		AllowHeaders:     c.AllowHeaders,
		AllowMethods:     c.AllowMethods,
	})
}

func NewCORS(v *viper.Viper) *CORSConfig {
	return &CORSConfig{
		AllowOrigins:     strings.Join(v.GetStringSlice("web.cors.allowed_origins"), ","),
		AllowMethods:     strings.Join(v.GetStringSlice("web.cors.allowed_methods"), ","),
		AllowHeaders:     strings.Join(v.GetStringSlice("web.cors.allowed_headers"), ","),
		AllowCredentials: v.GetBool("web.cors.allow_credentials"),
	}
}

func (c *CORSConfig) Register(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     c.AllowOrigins,
		AllowMethods:     c.AllowMethods,
		AllowHeaders:     c.AllowHeaders,
		AllowCredentials: c.AllowCredentials,
	}))
}
