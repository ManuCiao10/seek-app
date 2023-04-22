package controllers

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// return login.html page if user not logged in
func LoginGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user != nil {
			c.HTML(http.StatusBadRequest, "login.html",
				gin.H{
					"content": "Please logout first",
					"user":    user,
				})
			return
		}
		c.HTML(http.StatusOK, "login.html", gin.H{
			"content": "",
			"user":    user,
		})
	}
}

// Post request to login user and save session
func LoginPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user UserPostLogin

		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Println(user.Email, user.Password)
		// check user in database
		if user.Email == "dev" && user.Password == "dev" {
			log.Println("User is valid")

			sessionID := uuid.New().String()

			// Store session ID in Redis store
			session := sessions.Default(c)
			session.Set("sessionID", sessionID)
			session.Save()

			// Set session ID as a cookie
			c.SetCookie("sessionID", sessionID, 3600, "/", "", false, true)

			c.Redirect(http.StatusFound, "/")
		} else {
			log.Println("User is invalid")
			c.HTML(http.StatusBadRequest, "login.html",
				gin.H{
					"content": "Invalid email or password",
					"user":    user,
				})

		}

	}
}

// return index.html page
func IndexGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// session := sessions.Default(c)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "This is an index page...",
		})
	}
}
