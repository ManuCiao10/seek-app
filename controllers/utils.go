package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
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
	user UserPost // user from login form
	// snUser   UserPostSignup // user from signup form
	// AuthUser UserPostSignup // user from database AuthRequired middleware
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

func getClientIP(c *gin.Context) string {
	forwardHeader := c.Request.Header.Get("x-forwarded-for")
	firstAddress := strings.Split(forwardHeader, ",")[0]
	if net.ParseIP(strings.TrimSpace(firstAddress)) != nil {
		return firstAddress
	}
	return c.ClientIP()
}

func isIPRateLimited(ip string) bool {
	// Implement your rate-limiting logic here
	// You could use a database or cache to store the number of requests per IP address
	// For example, you could use Redis or Memcached to store a counter for each IP address
	// If the counter exceeds a certain threshold, you could return true to indicate that the IP is rate-limited
	// Otherwise, you could return false to indicate that the IP is not rate-limited
	return false
}
