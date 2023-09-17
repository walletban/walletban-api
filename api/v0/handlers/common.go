package handlers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"net/http"
	"walletban-api/api/v0/presenter"
)

func handleError(c *fiber.Ctx, err error, file string) error {
	errMessage := errors.New(fmt.Sprintf("%v: %v", file, err))
	log.Error(fmt.Sprintf("%v: %v", file, err))
	c.Status(http.StatusBadRequest)
	return c.JSON(presenter.Failure(errMessage))
}
