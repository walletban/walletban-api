package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"walletban-api/internal/entities"
	"walletban-api/internal/utils"
)

func main() {
	_ = connectToDb()
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "true",
			"details": "walletban-api",
			"author":  "Hemanth Krishna <@DarthBenro008>",
			"repository": "https://github." +
				"com/DarthBenro008/walletban",
		})
	})
	app.Listen(":8000")
}

func connectToDb() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Kolkata", utils.DBUrl, utils.DBUser, utils.DBPassword, utils.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&entities.User{}, &entities.Project{}, &entities.Consumer{})
	if err != nil {
		panic(err)
	}
	return db
}
