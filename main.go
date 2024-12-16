package main

import (
	"github.com/cenpnt/Todo-List-BackEnd/initializers"
	"github.com/cenpnt/Todo-List-BackEnd/routers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	r.GET("/users", routers.GetUsers)
	r.GET("/user/:id", routers.GetUserByID)
	r.POST("/signup", routers.SignUp)

	r.Run()
}