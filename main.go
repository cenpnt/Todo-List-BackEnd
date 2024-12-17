package main

import (
	"github.com/cenpnt/Todo-List-BackEnd/initializers"
	"github.com/cenpnt/Todo-List-BackEnd/middleware"
	"github.com/cenpnt/Todo-List-BackEnd/routers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	// User routers
	r.GET("/users", routers.GetUsers)
	r.GET("/user/:id", routers.GetUserByID)
	r.POST("/signup", routers.SignUp)
	r.POST("/login", routers.Login)
	
	// Task routers
	r.POST("/task", middleware.AuthMiddleware, routers.CreateTask)

	r.Run()
}