package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/st107853/fast_reading/controllers"
	"github.com/st107853/fast_reading/middleware"
	"github.com/st107853/fast_reading/services"
)

type BookRouteController struct {
	bookController controllers.BookController
}

func NewBookRouteController(bookController controllers.BookController) BookRouteController {
	return BookRouteController{bookController}
}

func (bc *BookRouteController) BookRoute(rg *gin.RouterGroup, bookService services.BookService, userService services.UserService) {
	rg.Use(middleware.DeserializeUser(userService))
	rg.POST("/", bc.bookController.CreateBook)
	rg.PUT("/:book_id", bc.bookController.UpdateBook)
	rg.PUT("/:book_id/:chapter_id/:last_index", bc.bookController.BookMark)
	rg.PUT("/release/:book_id", bc.bookController.ReleaseBook)
	rg.DELETE("/:book_id", bc.bookController.DeleteBook)
	rg.DELETE("/", bc.bookController.DeleteAllBooks)
	rg.DELETE("/chapter/:chapter_id", bc.bookController.DeleteChapter)
	rg.GET("/", bc.bookController.AllBooks)
	rg.GET("/continue", bc.bookController.ContinueReading)
	rg.GET("/book/:book_id", bc.bookController.GetBook)
	rg.GET("/book/:book_id/:chapter_id/:last_index", bc.bookController.GetChapter)
	rg.POST("/book/:book_id/favourite", bc.bookController.BookFavourite)
	rg.GET("/addbook", bc.bookController.AddBook)
	rg.GET("/addbook/:book_id", bc.bookController.EditBook)
	rg.GET("/addbook/:book_id/chapter", bc.bookController.AddBookChapter)
	rg.POST("/addbook/:book_id/chapter", bc.bookController.CreateChapter)
	rg.GET("/addbook/:book_id/chapter/:chapter_id", bc.bookController.EditBookChapter)
	rg.PUT("/addbook/:book_id/chapter/:chapter_id", bc.bookController.UpdateBookChapter)
	rg.PUT("/book/:book_id/labels", bc.bookController.AddLabel)
	rg.GET("/filter/", bc.bookController.ListAllBooks)
}
