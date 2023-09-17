package middlewares

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"walletban-api/api/v0/presenter"
	"walletban-api/internal/utils"
)

func JwtMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(utils.JwtSecret),
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			ctx.Status(http.StatusUnauthorized)
			return ctx.JSON(presenter.Failure(errors.New("invalid token")))
		},
		SuccessHandler: func(ctx *fiber.Ctx) error {
			user := ctx.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			username := claims["username"].(string)
			uid := claims["uid"].(float64)
			pid := claims["pid"].(float64)
			ctx.Locals("username", username)
			ctx.Locals("uid", uid)
			ctx.Locals("pid", pid)
			return ctx.Next()
		},
	})
}
