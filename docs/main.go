package main

import (
	"log"
	"os"

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
		log.Fatal("Error loading .env file")
	}

	appConfig.AppPort = getEnv("APP_PORT", "9000")

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	public := router.Group("/")
	routes.PublicRoutes(public)

	private := router.Group("/")
	private.Use(middleware.AuthRequired)
	routes.PrivateRoutes(private)

	database.InitMongoDB()

	log.Fatal(router.Run(":" + appConfig.AppPort))

	log.Println("Server started on port " + appConfig.AppPort)

}

//add design part CSS
//add login discord and google
//TODO create 2 distingushed index.html (one for the user not logged in and one for the user logged in)

//#########URLS##########
//https://www.youtube.com/watch?v=7hOfR6wHMaw
