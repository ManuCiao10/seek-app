package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/seek/docs/controllers"
	"github.com/seek/docs/database"
	"github.com/seek/docs/middleware"
	routes "github.com/seek/docs/routers"
)

var appConfig = controllers.AppConfig{}

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

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	store, err := redis.NewStore(10, "tcp", "localhost:6379", "", []byte(appConfig.AppSecret))
	if err != nil {
		panic(err)
	}

	router.Use(sessions.Sessions("session", store))

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	router.Use(sessions.Sessions("session", cookie.NewStore([]byte(appConfig.AppSecret))))

	public := router.Group("/")
	routes.PublicRoutes(public)

	private := router.Group("/")
	private.Use(middleware.AuthRequired)
	routes.PrivateRoutes(private)

	database.InitMongoDB()

	router.Run(":" + appConfig.AppPort)

	log.Println("Server started on port " + appConfig.AppPort)

}

//improve control user logged in
//add design part CSS
//different index bewtwenn user logged in and not logged in (done)
//add login discord and google
//export GIN_MODE=release
//check how to assign cookies and sessions
