package handlers

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"walletban-api/api/v0/presenter"
	"walletban-api/internal/entities"
	"walletban-api/internal/services"
	"walletban-api/internal/utils"
)

const (
	projectErr = "project error"
)

func CreateProject(ctx context.Context, service services.ApplicationService, username string, userID uint) (uint, error) {
	clientId, err := utils.GenerateRandomString(utils.ClientIDRandomLength)
	if err != nil {
		return 0, err
	}
	clientSecret, err := utils.GenerateRandomString(utils.ClientSecretRandomLength)
	if err != nil {
		return 0, err
	}
	apiKey, err := utils.GenerateRandomString(6)
	if err != nil {
		return 0, err
	}
	project := entities.Project{
		UserID:       userID,
		Name:         username + "'s Project",
		TokenName:    "",
		Consumers:    nil,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		ApiKey:       apiKey,
	}
	res, err := service.ProjectRepository.Create(ctx, project)
	if err != nil {
		return 0, err
	}
	return res.ID, nil
}

func GetDashboard(service services.ApplicationService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		username := fmt.Sprintf("%v", ctx.Locals("username"))
		uid := uint(int(ctx.Locals("uid").(float64)))
		pid := uint(int(ctx.Locals("pid").(float64)))
		user := entities.User{Username: username}
		user.ID = uid
		userData, err := service.UserRepository.FindOne(ctx.Context(), user)
		if err != nil {
			return handleError(ctx, err, "dashboard error")
		}
		project := entities.Project{UserID: uid}
		project.ID = pid
		projectData, err := service.ProjectRepository.FindOne(ctx.Context(), project)
		if err != nil {
			return handleError(ctx, err, "dashboard error")
		}
		userData.Project = *projectData
		return ctx.JSON(presenter.Success(userData, "user fetched!"))
	}
}

//TODO: List all Consumers and ProjectData
//TODO: block and unblock users
//TODO: Generate new apikey
