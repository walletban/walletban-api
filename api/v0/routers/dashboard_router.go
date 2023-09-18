package routers

import (
	"github.com/gofiber/fiber/v2"
	"walletban-api/api/v0/handlers"
	"walletban-api/internal/services"
)

func DashboardRouter(app fiber.Router, applicationService services.ApplicationService) {
	app.Get("/dashboard", handlers.GetDashboard(applicationService))
}
