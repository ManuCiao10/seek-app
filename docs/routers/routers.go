package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/seek/docs/controllers"
)

// User
func UserRouters(incomingRouters *gin.Engine) {
	incomingRouters.GET("/users", controllers.GetUsers)
	// incomingRouters.GET("/users/:id", controllers.GetUser)
	// incomingRouters.POST("/users", controllers.AddUser)
	// incomingRouters.PUT("/users/:id", controllers.UpdateUser)
	// incomingRouters.DELETE("/users/:id", controllers.DeleteUser)

}
