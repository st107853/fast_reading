package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/st107853/fast_reading/controllers"
	"github.com/st107853/fast_reading/middleware"
	"github.com/st107853/fast_reading/services"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup, userService services.UserService) {

	router := rg.Group("users")
	router.Use(middleware.DeserializeUser(userService))
	router.GET("/me", uc.userController.GetMe)
}
