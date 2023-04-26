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

	//Error: [sessions] ERROR! securecookie: the value is not valid we should use a static key
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
//TODO improve database client implementation

// https://www.youtube.com/watch?v=7hOfR6wHMaw
// https://github.com/Skarlso/google-oauth-go-sample/blob/master/database/mongo.go
// https://skarlso.github.io/2016/06/12/google-signin-with-go/
// https://skarlso.github.io/2016/11/02/google-signin-with-go-part2/
// https://github.com/zalando/gin-oauth2/blob/47b9fc0cb1395111098062ff8d991174fa40f6b3/google/google.go#L99
