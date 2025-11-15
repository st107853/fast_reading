package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/st107853/fast_reading/controllers"
	"github.com/st107853/fast_reading/services"
)

type BookRouteController struct {
	bookController controllers.BookController
}

func NewBookRouteController(bookController controllers.BookController) BookRouteController {
	return BookRouteController{bookController}
}

func (bc *BookRouteController) BookRoute(rg *gin.RouterGroup, bookService services.BookService) {
	rg.POST("/", bc.bookController.CreateBook)
	rg.PUT("/:id", bc.bookController.UpdateBook)
	rg.DELETE("/:id", bc.bookController.DeleteBook)
	rg.DELETE("/", bc.bookController.DeleteAllBooks)
	rg.GET("/", bc.bookController.AllBooks)
	rg.GET("/one/:book_id", bc.bookController.GetBook)
	rg.GET("/one/chapter/:chapter_id", bc.bookController.GetChapter)
	rg.POST("/one/:book_id/favourite", bc.bookController.BookFavourite)
	rg.GET("/all/:name", bc.bookController.GetBooksByName)
	rg.GET("/addbook", bc.bookController.AddBook)
	rg.GET("/one/:book_id/chapter", bc.bookController.AddBookChapter)
	rg.POST("/one/:book_id/chapter", bc.bookController.CreateChapter)
	rg.GET("/edit/:book_id", bc.bookController.EditBook)
}
