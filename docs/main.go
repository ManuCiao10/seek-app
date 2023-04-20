package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/seek/docs/controllers"
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
		fmt.Println(err)
	}

	appConfig.AppName = getEnv("APP_NAME", "Seek")
	appConfig.AppPort = getEnv("APP_PORT", "9000")
	appConfig.AppURL = getEnv("APP_URL", "http://localhost:9000")

	router := gin.New()
	router.Use(gin.Logger())
	// gin.SetMode(gin.ReleaseMode)

	router.GET("/", controllers.HomePage)
	// router.GET("/login", controllers.Login)
	// router.GET("/register", controllers.Register)
	// router.GET("/logout", controllers.Logout)
	// router.GET("/profile", controllers.Profile)
	// router.GET("/profile/edit", controllers.ProfileEdit)
	// router.GET("/profile/password", controllers.ProfilePassword)
	// router.GET("/profile/delete", controllers.ProfileDelete)
	// router.GET("/profile/verify", controllers.ProfileVerify)

	router.Run(":" + appConfig.AppPort)

}

//add login discord and google
//add logTail for log
//export GIN_MODE=release
