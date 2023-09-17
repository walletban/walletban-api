package utils

import "os"

var (
	DBUrl                    = os.Getenv("DB_SERVER")
	DBUser                   = os.Getenv("DB_USER")
	DBPassword               = os.Getenv("DB_PASSWORD")
	DBName                   = "walletban"
	JwtSecret                = os.Getenv("JWT_SECRET")
	RedirectUrl              = os.Getenv("REDIRECT_URL")
	FrontendUrl              = "https://walletban.xyz"
	ClientIDRandomLength     = 20
	ClientSecretRandomLength = 10
	OauthBypass              = "glfyfe"
)
