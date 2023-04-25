package main

import (
	"log"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/seek/docs/controllers"
	"github.com/seek/docs/database"
	"github.com/seek/docs/middleware"
	routes "github.com/seek/docs/routers"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	token, err := controllers.RandToken(64)
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

//TODO display all the error message in the html and the success message (like account created)
//TODO login discord and google
//TODO fixing the redirection html (user not logged in and user logged in)

//#########URLS##########
//https://www.youtube.com/watch?v=7hOfR6wHMaw
