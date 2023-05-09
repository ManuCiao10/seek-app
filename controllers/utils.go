package controllers

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/oauth2"
)

var User_google UserGoogle
var confgoogle *oauth2.Config

var User_discord UserDiscord
var confdiscord *oauth2.Config

// Endpoint is Discord's OAuth 2.0 endpoint.
var Endpoint = oauth2.Endpoint{
	AuthURL:   "https://discord.com/api/oauth2/authorize",
	TokenURL:  "https://discord.com/api/oauth2/token",
	AuthStyle: oauth2.AuthStyleInParams,
}

var (
	user     UserPostLogin  // user from login form
	dbUser   UserPostLogin  // user from database
	snUser   UserPostSignup // user from signup form
	AuthUser UserPostSignup // user from database AuthRequired middleware
)

// RandToken generates a random @l length token.
func RandToken(l int) (string, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func getLoginURL(state string) string {
	return confgoogle.AuthCodeURL(state)
}

func getDiscordLoginURL(state string) string {
	return confdiscord.AuthCodeURL(state)
}
