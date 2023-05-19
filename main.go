package main

import (
	"log"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/seek/controllers"
	"github.com/seek/database"
	"github.com/seek/middleware"
	routes "github.com/seek/routers"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	token, err := controllers.RandToken(32)
	if err != nil {
		log.Fatal("unable to generate random token: ", err)
	}

	store := sessions.NewCookieStore([]byte(token))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})

	router.Use(sessions.Sessions("session", store))

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.Static("/assets", "./assets")

	public := router.Group("/")
	routes.PublicRoutes(public)

	private := router.Group("/")
	private.Use(middleware.AuthRequired)
	routes.PrivateRoutes(private)

	database.InitMongoDB()

	log.Fatal(router.Run(":9000"))
}
