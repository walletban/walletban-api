package routers

import (
	"github.com/gofiber/fiber/v2"
	"walletban-api/api/v0/handlers"
	"walletban-api/internal/services"
)

func ConsumerRouter(app fiber.Router, applicationService services.ApplicationService) {
	app.Post("/consumer", handlers.RegisterConsumer(applicationService))
	app.Post("/invoke", handlers.InvokeContract(applicationService))
}
