package main

import (
	"log"

	"github.com/st107853/fast_reading/controllers"
	"github.com/st107853/fast_reading/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Static("/static", "./")

	r.Use(gin.Logger())

	err := models.Connect()
	if err != nil {
		log.Fatal(err)
	}

	library := r.Group("/library")
	{
		library.POST("/", controllers.CreateBook)
		library.PUT("/:id", controllers.UpdateBook)
		library.DELETE("/:id", controllers.DeleteBook)
		library.DELETE("/", controllers.DeleteAllBooks)
		library.GET("/", controllers.AllBooks)
		library.GET("/one/:id", controllers.GetBook)
		library.GET("/all/:name", controllers.GetBooks)
		// library.POST("/multiple", controllers.InsertMultipleBooks)
	}

	log.Println("Server started")
	r.Run()
}
