package middleware

import (
	"github.com/gin-gonic/gin"
)

func AuthRequired(c *gin.Context) {
	// session := sessions.Default(c)
	// user := session.Get("user")
	// if user == nil {
	// 	log.Println("User not logged in")
	// 	c.Redirect(http.StatusMovedPermanently, "/login")
	// 	c.Abort()
	// 	return
	// }
	c.Next()
}
