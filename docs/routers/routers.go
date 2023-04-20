package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seek/docs/controllers"
)

func PublicRoutes(g *gin.RouterGroup) {

	g.GET("/login", controllers.LoginGetHandler())
	// g.POST("/login", controllers.LoginPostHandler())
	// g.GET("/register", controllers.RegisterGetHandler())
	// g.POST("/register", controllers.RegisterPostHandler())
	g.GET("/", controllers.IndexGetHandler())

}

func PrivateRoutes(g *gin.RouterGroup) {

	// g.GET("/dashboard", controllers.DashboardGetHandler())
	// g.GET("/logout", controllers.LogoutGetHandler())

}
