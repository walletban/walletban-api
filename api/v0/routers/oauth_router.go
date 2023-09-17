package routers

import (
	"github.com/gofiber/fiber/v2"
	"walletban-api/api/v0/handlers"
	"walletban-api/internal/services"
)

func OAuthRouter(app fiber.Router, applicationService services.ApplicationService) {
	app.Get("/auth/google/login", handlers.GoogleLogin())
	app.Get("/auth/google/callback", handlers.GoogleCallback(applicationService))
}
