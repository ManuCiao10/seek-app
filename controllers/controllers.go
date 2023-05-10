package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seek/database"
)

/*
determine if the user is logged in or not by searching the sessionID in the cookie and the database MongoDB
user logged in ==> display index.html with [sign up/log in button] + [sell now button]
user not logged in ==> display index.html with [profile picture] + [sell now button] + [messages]
*/
func IndexGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("IndexGetHandler: %v", c.Request.URL.Path)

		sessionID, err := c.Cookie("sessionID")

		if len(sessionID) == 0 {
			log.Printf("User Token header is missing, user not logged in")

			c.HTML(http.StatusOK, "not_logged_in.html", gin.H{
				"content": "",
			})
			return
		}

		if err != nil {
			log.Printf("Error IndexGetHandler gettin sessionID: user is not logged in")

			c.HTML(http.StatusOK, "not_logged_in.html", gin.H{
				"content": "",
			})
			return
		}

		log.Printf("Found session ID [%v]", sessionID)
		err = database.CheckSession(c, sessionID)

		if err != nil {
			c.HTML(http.StatusOK, "not_logged_in.html", gin.H{})
			return
		}

		err = database.CheckSessionExpired(c)

		if err != nil {
			c.HTML(http.StatusOK, "not_logged_in.html", gin.H{})
			return
		}

		log.Printf("User logged in, session ID: %v", sessionID)

		c.HTML(http.StatusOK, "logged_in.html", gin.H{})

	}
}
