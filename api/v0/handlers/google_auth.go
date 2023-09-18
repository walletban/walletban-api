package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"os"
	"strconv"
	"strings"
	"walletban-api/api/v0/presenter"
	"walletban-api/internal/entities"
	"walletban-api/internal/services"
	"walletban-api/internal/utils"
)

var (
	state = "holderState"
)

const (
	oauthErr = "oauth error"
)

type googleAuthResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func oAuthGoogleConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  utils.RedirectUrl + "/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleLogin(service services.ApplicationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		redirect := c.Query("redirect", utils.RedirectUrl)
		clientID := c.Query("clientId", "")
		clientSecret := c.Query("clientSecret", "")
		if clientID == "" || clientSecret == "" {
			return handleError(c, errors.New("invalid client id/secret"), "Invalid Client ID/Secret")
		}
		if clientID != utils.OauthBypass {
			proj := entities.Project{ClientSecret: clientSecret, ClientId: clientID}
			res, err := service.ProjectRepository.FindOne(c.Context(), proj)
			if err != nil {
				return handleError(c, err, "Client ID/Secret not found")
			}
			if res == nil {
				return handleError(c, errors.New("invalid client id/secret"), "Client ID/Secret not found")
			}
		}
		tempState, err := utils.GenerateOauthHash(redirect + ";" + clientID + ";" + clientSecret)
		state = tempState
		if err != nil {
			return handleError(c, err, "unknown error has occurred")
		}
		url := oAuthGoogleConfig().AuthCodeURL(state)
		return c.Redirect(url, http.StatusTemporaryRedirect)
	}
}

func GoogleCallback(service services.ApplicationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		redirectData, err := utils.ValidateToken(c.FormValue("state", ""))
		splitRedirectData := strings.Split(redirectData, ";")
		url := splitRedirectData[0]
		clientId := splitRedirectData[1]
		clientSecret := splitRedirectData[2]
		log.Info(fmt.Sprintf("Requested for URL redirect at: %s", url))
		if err != nil {
			return c.Redirect(utils.FrontendUrl, http.StatusTemporaryRedirect)
		}
		token, err := oAuthGoogleConfig().Exchange(context.Background(), c.FormValue("code"))
		if err != nil {
			fmt.Print(err)
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.Failure(err))
		}

		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			return c.SendString("Cannot get your details bro")
		}
		defer resp.Body.Close()
		googleResponse := googleAuthResponse{}
		err = json.NewDecoder(resp.Body).Decode(&googleResponse)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.Failure(err))
		}
		// Check if user login or consumer login
		if clientId == utils.OauthBypass {
			return userLogin(c, service, googleResponse, url)
		} else {
			// Consumer Login
			return consumerLogin(c, service, googleResponse, url, clientSecret, clientId)
		}
	}
}

func consumerLogin(c *fiber.Ctx, service services.ApplicationService, googleResponse googleAuthResponse, url string, clientSecret string, clientId string) error {
	proj := entities.Project{ClientSecret: clientSecret, ClientId: clientId}
	projData, err := service.ProjectRepository.FindOne(c.Context(), proj)
	if err != nil {
		return handleError(c, err, "Client ID/Secret not found")
	}
	consumer := entities.Consumer{
		ProjectID:           projData.ID,
		Name:                googleResponse.GivenName + " " + googleResponse.FamilyName,
		Email:               googleResponse.Email,
		WalletGKey:          "",
		WalletEncryptedSKey: "",
		PfpUrl:              googleResponse.Picture,
	}
	res, err := service.ConsumerRepository.Create(c.Context(), consumer)
	if err != nil {
		return handleError(c, err, "unable to create consumer")
	}
	unique := url + "?token=" + strconv.Itoa(int(res.ID)) + "&first=" + strconv.FormatBool(res.IsFirstTime) + "&isActivated=" + strconv.FormatBool(res.IsWalletActivated)
	return c.Redirect(unique, http.StatusTemporaryRedirect)
}

func userLogin(c *fiber.Ctx, service services.ApplicationService, googleResponse googleAuthResponse, url string) error {
	var userData = entities.User{
		Name:     googleResponse.GivenName + " " + googleResponse.FamilyName,
		Username: strings.Split(googleResponse.Email, "@")[0],
		PfpUrl:   googleResponse.Picture,
		Email:    googleResponse.Email,
	}
	userInsertion, err := service.UserRepository.Create(c.Context(), userData)
	if err != nil {
		var registeredUser entities.User
		registeredUser.Username = userData.Username
		userInsertion, err := service.UserRepository.FindOne(c.Context(), registeredUser)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.Failure(err))
		}
		if userInsertion.Project.ID <= uint(0) {
			id, err := CreateProject(c.Context(), service, userInsertion.Name, userInsertion.ID)
			if err != nil {
				fmt.Println(err)
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenter.Failure(err))
			}
			userInsertion.Project.ID = id
		}
		jwtToken := userInsertion.GetSignedJWT()
		unique := url + "?token=" + jwtToken + "&first=" + strconv.FormatBool(userInsertion.IsFirstTime)
		return c.Redirect(unique, http.StatusTemporaryRedirect)
	}
	id, err := CreateProject(c.Context(), service, userInsertion.Name, userInsertion.ID)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return c.JSON(presenter.Failure(err))
	}
	userInsertion.Project.ID = id
	jwtToken := userInsertion.GetSignedJWT()
	unique := url + "?token=" + jwtToken + "&first=" + strconv.FormatBool(userInsertion.IsFirstTime)
	return c.Redirect(unique, http.StatusTemporaryRedirect)
}
