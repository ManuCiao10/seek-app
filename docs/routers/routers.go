package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seek/docs/controllers"
	"github.com/seek/pkg/auth"
)

func PublicRoutes(g *gin.RouterGroup) {

	g.GET("/", controllers.IndexGetHandler())
	g.GET("/login", controllers.LoginGetHandler())
	g.POST("/login", controllers.LoginPostHandler())

	g.GET("/login/google", auth.HandleGoogleLogin())
	g.GET("/login/callback-google", controllers.LoginGetHandler())

	// g.GET("/register", controllers.RegisterGetHandler())
	// g.POST("/register", controllers.RegisterPostHandler())
	// g.GET("/member/:id", controllers.MemberGetHandler())

}

func PrivateRoutes(g *gin.RouterGroup) {

	g.GET("/dashboard", controllers.DashboardGetHandler())
	// g.GET("/logout", controllers.LogoutGetHandler())
	//inbox
	// g.GET("/inbox/:id", controllers.InboxGetHandler())

}
