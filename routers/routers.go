package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seek/controllers"
)

// Display all routes
func PublicRoutes(g *gin.RouterGroup) {

	g.GET("/", controllers.IndexGetHandler())
	g.GET("/login", controllers.LoginGetHandler())
	g.POST("/login", controllers.LoginPostHandler())
	g.GET("/signup", controllers.SignupGetHandler())
	g.POST("/signup", controllers.SignupPostHandler())

	g.GET("/login/google", controllers.HandleGoogleLogin())
	g.GET("/login/google-callback", controllers.HandleGoogleCallback())

	g.GET("/login/discord", controllers.HandleDiscordLogin())
	g.GET("/login/callback-discord", controllers.HandleDiscordCallback())

	// g.GET("/login/apple", controllers.LoginGetHandler())
	// g.GET("/login/callback-apple", controllers.LoginGetHandler())

}

// Diplay only if user is logged in (middleware)
func PrivateRoutes(g *gin.RouterGroup) {

	// g.GET("/dashboard", controllers.DashboardGetHandler())
	// g.GET("/logout", controllers.LogoutGetHandler())
	// g.GET("/inbox", controllers.InboxGetHandler())
	// g.GET("/inbox/:id", controllers.InboxGetHandler())

	// POST /products – to add a new product
	// PUT /products/:productId – to update a product
	// DELETE /products/:productId – to delete a product

	// g.GET("/logout", controllers.LogoutGetHandler())

}
