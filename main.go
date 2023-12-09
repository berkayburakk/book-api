package main

import (
	"book-api/application/controllers"
	_ "book-api/docs"
	"book-api/initiliazers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initiliazers.LoadEnvVariables()
	initiliazers.ConnectToDB()
}

// @title Book Service API
// @version 1.0
// @description A book service API in Go Using Gin Framework

// @host localhost:3000
// @BasePath /
func main() {
	r := gin.Default()
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/books", controllers.CreateBook)
	r.GET("/books", controllers.GetBooks)
	r.GET("/books/:id", controllers.GetBookById)
	r.PUT("/books/:id", controllers.UpdateBook)
	r.DELETE("/books/:id", controllers.DeleteBook)

	r.Run()
}
