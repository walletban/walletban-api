package handlers

import (
	"context"
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
