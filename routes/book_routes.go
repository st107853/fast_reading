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
	rg.PUT("/:book_id", bc.bookController.UpdateBook)
	rg.DELETE("/:book_id", bc.bookController.DeleteBook)
	rg.DELETE("/", bc.bookController.DeleteAllBooks)
	rg.DELETE("/chapter/:chapter_id", bc.bookController.DeleteChapter)
	rg.GET("/", bc.bookController.AllBooks)
	rg.GET("/book/:book_id", bc.bookController.GetBook)
	rg.GET("/book/:book_id/:chapter_id", bc.bookController.GetChapter)
	rg.POST("/book/:book_id/favourite", bc.bookController.BookFavourite)
	rg.GET("/all/:name", bc.bookController.GetBooksByName)
	rg.GET("/addbook", bc.bookController.AddBook)
	rg.GET("/addbook/:book_id", bc.bookController.EditBook)
	rg.GET("/addbook/:book_id/chapter", bc.bookController.AddBookChapter)
	rg.POST("/addbook/:book_id/chapter", bc.bookController.CreateChapter)
	rg.GET("/addbook/:book_id/chapter/:chapter_id", bc.bookController.EditBookChapter)
	rg.PUT("/addbook/:book_id/chapter/:chapter_id", bc.bookController.UpdateBookChapter)
	rg.POST("/addbook/:book_id/label/:label_id", bc.bookController.AddLabel)
	rg.GET("/filter/", bc.bookController.ListAllBooks)
}
