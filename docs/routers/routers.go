package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seek/docs/controllers"
)

func PublicRoutes(g *gin.RouterGroup) {

	g.GET("/", controllers.IndexGetHandler())
	g.GET("/login", controllers.LoginGetHandler())
	g.POST("/login", controllers.LoginPostHandler())

	// g.GET("/login/google", auth.HandleGoogleLogin())
	// g.GET("/login/callback-google", controllers.LoginGetHandler())

	// g.GET("/register", controllers.RegisterGetHandler())
	// g.POST("/register", controllers.RegisterPostHandler())
	// g.GET("/member/:id", controllers.MemberGetHandler())
	// GET /products – to list several products
	// GET /products/:productId – to get details of one product

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

}
