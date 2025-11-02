package controllers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/st107853/fast_reading/models"
	"github.com/st107853/fast_reading/services"
)

var userPage = template.Must(template.New("user_page.html").Funcs(template.FuncMap{
	"extractNumericPart": extractNumericPart,
}).ParseFiles("./static/user_page.html"))

type UserController struct {
	userService services.UserService
}

type UserData struct {
	Name      string
	Favourite []*models.Book
}

func NewUserController(userService services.UserService) UserController {
	return UserController{userService}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.User)

	//var user = models.FilteredResponse(currentUser)

	data := UserData{
		Name:      currentUser.Name,
		Favourite: currentUser.FavoriteBooks,
	}

	// Execute the template and write the output to the response writer
	if err := userPage.Execute(ctx.Writer, data); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": models.FilteredResponse(currentUser)}})
}
