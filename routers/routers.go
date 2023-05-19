package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seek/controllers"
)

func PublicRoutes(g *gin.RouterGroup) {
	g.GET("/", controllers.IndexGetHandler())

	g.GET("/login", controllers.LoginGetHandler())
	g.POST("/login", controllers.LoginPostHandler())

	g.GET("/signup", controllers.SignupGetHandler())
	g.POST("/signup", controllers.SignupPostHandler())

	g.GET("/login/google", controllers.HandleGoogleLogin())
	g.GET("/login/google-callback", controllers.HandleGoogleCallback())

	g.GET("/login/discord", controllers.HandleDiscordLogin())
	g.GET("/login/discord-callback", controllers.HandleDiscordCallback())

}

func PrivateRoutes(g *gin.RouterGroup) {

	g.GET("/logout", controllers.LogoutGetHandler())
	// g.GET("/dashboard", controllers.DashboardGetHandler())
	// g.GET("/inbox", controllers.InboxGetHandler())
	// g.GET("/inbox/:id", controllers.InboxGetHandler())

	// POST /products – to add a new product
	// PUT /products/:productId – to update a product
	// DELETE /products/:productId – to delete a product

	// g.GET("/settings", controllers.SettingsGetHandler())

	// g.GET("/settings/profile", controllers.SettingsProfileGetHandler())
	// g.POST("/settings/profile", controllers.SettingsProfilePostHandler())

	// g.GET("/settings/account", controllers.SettingsAccountGetHandler())
	// g.POST("/settings/account", controllers.SettingsAccountPostHandler())

	// g.GET("/settings/security", controllers.SettingsSecurityGetHandler())
	// g.POST("/settings/security", controllers.SettingsSecurityPostHandler())

	// g.GET("/logout", controllers.LogoutGetHandler())

	// g.GET("/api/v2/catalogs", controllers.CatalogsGetHandler())
	// g.GET("/api/v2/catalogs/:id", controllers.CatalogGetHandler())
	// g.GET("/api/v2/users/:id", controllers.UserGetHandler())

}
