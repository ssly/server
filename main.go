package main

import (
	"./controllers"
	"./middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()

	route.Use(middleware.ConnectDB)

	route.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome golang")
	})

	route.POST("/task/manager", controllers.CreateTask)
	route.DELETE("/task/manager/:id", controllers.DeleteTask)
	route.DELETE("/task/manager", controllers.DeleteTask)
	route.PUT("/task/manager/:id", controllers.UpdateTask)
	route.PUT("/task/manager", controllers.UpdateTask)
	route.GET("/task/manager/:id", controllers.GetTask) // find one task
	route.GET("/task/manager", controllers.GetTask)     // find all task list

	route.Run(":5432")
}
