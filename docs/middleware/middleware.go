package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	isLoggedIn := session.Get("isLoggedIn")

	if isLoggedIn != true {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	} else {
		c.Next()
	}
	c.Next()
}
