package entities

import "os"

var (
	SecretKey  = os.Getenv("SECRET_KEY")
	CookieName = "USER_COOKIE"
)

type Login struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
