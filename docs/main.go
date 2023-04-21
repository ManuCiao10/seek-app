package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/seek/docs/controllers"
	"github.com/seek/docs/middleware"
	routes "github.com/seek/docs/routers"
)

var (
	appConfig = controllers.AppConfig{}
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	appConfig.AppName = getEnv("APP_NAME", "Seek")
	appConfig.AppPort = getEnv("APP_PORT", "9000")
	appConfig.AppURL = getEnv("APP_URL", "http://localhost:9000")
	appConfig.AppSecret = getEnv("APP_SECRET", "seek")

	router := gin.New()
	router.Use(gin.Logger())

	//load template && static
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	//save cookie session
	router.Use(sessions.Sessions("session", cookie.NewStore([]byte(appConfig.AppSecret))))

	//load routes
	public := router.Group("/")
	routes.PublicRoutes(public)

	//load middleware
	private := router.Group("/")
	private.Use(middleware.AuthRequired)
	routes.PrivateRoutes(private)

	router.Run(":" + appConfig.AppPort)

}

//add login discord and google
//add logTail for log
//export GIN_MODE=release
