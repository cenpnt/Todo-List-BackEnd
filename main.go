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
	r.GET("/tasks", middleware.AuthMiddleware, routers.GetAllTasks)
	r.POST("/create-task", middleware.AuthMiddleware, routers.CreateTask)
	r.PATCH("/toggleTask/:id", middleware.AuthMiddleware, routers.ToggleTask)
	r.PUT("/edit-task/:id", middleware.AuthMiddleware, routers.EditTask)
	r.DELETE("/delete-task/:id", middleware.AuthMiddleware, routers.DeleteTask)

	r.Run()
}