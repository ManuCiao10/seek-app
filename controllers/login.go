package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/seek/database"
)

// return login.html page if user not logged in
func LoginGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("LoginGetHandler: %v", c.Request.URL.Path)
		sessionID, err := c.Cookie("sessionID")

		if err != nil {
			log.Printf("User is not logged in, redirect to login page")

			c.HTML(http.StatusOK, "login.html", gin.H{})
			return
		}

		log.Printf("Found session ID in cookie %v", sessionID)

		err = database.CheckSession(c, sessionID)
		if err != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{})

			return
		}

		err = database.CheckSessionExpired(c)
		if err != nil {
			if err != nil {
				c.HTML(http.StatusOK, "login.html", gin.H{})
				return
			}
		}

		log.Printf("User is logged in, redirect to index")
		log.Printf("SessionID: %v", sessionID)

		c.HTML(http.StatusOK, "logged_in.html", gin.H{})

	}
}

// Post request to login user:
// if user is valid (email and password) => store session ID in MongoDB and as a cookie
// else return login.html page with error message
func LoginPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		log.Printf("LoginPostHandler: %v", c.Request.URL.Path)

		email := c.PostForm("email")
		password := c.PostForm("password")

		// ip := getClientIP(c)

		// if isIPRateLimited(c.ClientIP()) {
		// 	c.JSON(http.StatusTooManyRequests, gin.H{"message": "Too many requests"})
		// 	return
		// }

		err := database.CheckEmail(c, email)

		if err != nil {
			log.Printf("Email invalid: %v", err)

			c.HTML(http.StatusBadRequest, "login.html", gin.H{})
			return
		}

		err = database.CheckPassword(c, password)

		if err != nil {
			log.Printf("Password invalid: %v", err)

			c.HTML(http.StatusBadRequest, "login.html", gin.H{})
			return
		}

		sessionID := uuid.NewString()
		expiresAt := time.Now().Add(15 * 24 * time.Hour)

		err = database.StoreSession(c, sessionID, expiresAt, email)
		if err != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{})
			return
		}

		c.SetCookie("sessionID", sessionID, int(expiresAt.Unix()), "/", "", false, true)

		log.Printf("Session ID: %v %v", sessionID, expiresAt)
		c.Redirect(http.StatusFound, "/")

	}
}
