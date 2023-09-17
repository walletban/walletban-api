package handlers

import (
	"context"
	"encoding/json"
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

func GoogleLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		redirect := c.Query("redirect", utils.RedirectUrl)
		fmt.Println(redirect)
		tempState, err := utils.GenerateOauthHash(redirect)
		state = tempState
		if err != nil {
			fmt.Println(err)
			return c.SendString("Some error has occurred.")
		}
		url := oAuthGoogleConfig().AuthCodeURL(state)
		return c.Redirect(url, http.StatusTemporaryRedirect)
	}
}

func GoogleCallback(service services.ApplicationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		url, err := utils.ValidateToken(c.FormValue("state", ""))
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
			jwtToken := userInsertion.GetSignedJWT()
			unique := url + "?token=" + jwtToken + "&first=" + strconv.FormatBool(userInsertion.IsFirstTime)
			return c.Redirect(unique, http.StatusTemporaryRedirect)
		}
		jwtToken := userInsertion.GetSignedJWT()
		unique := url + "?token=" + jwtToken + "&first=" + strconv.FormatBool(userInsertion.IsFirstTime)
		return c.Redirect(unique, http.StatusTemporaryRedirect)
	}
}
