package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/go-redis/redis"
	"github.com/st107853/fast_reading/config"
	"github.com/st107853/fast_reading/controllers"
	"github.com/st107853/fast_reading/models"
	"github.com/st107853/fast_reading/routes"
	"github.com/st107853/fast_reading/services"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	redisclient *redis.Client

	userService         services.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authCollection      *mongo.Collection
	authService         services.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	booksCollection     *mongo.Collection
	bookService         services.BookService
	BookController      controllers.BookController
	BookRouteController routes.BookRouteController
)

func init() {
	gin.SetMode(gin.ReleaseMode) // Set Gin to release mode for production

	ctx = context.TODO()

	err := models.ConnectToMongoDB()
	if err != nil {
		log.Fatal("Could not connect to MongoDB", err)
	}

	err = models.ConnectToRedis() // Assign Redis client to redisclient
	if err != nil {
		log.Fatal("Could not connect to Redis")
	}

	// Collections
	authCollection = models.DB.Database("golang_mongodb").Collection("users")
	booksCollection = models.DB.Database("golang_mongodb").Collection("books")
	userService = services.NewUserServiceImpl(authCollection, ctx)
	authService = services.NewAuthService(authCollection, ctx)
	bookService = services.NewBookService(booksCollection, ctx)

	AuthController = controllers.NewAuthController(authService, userService)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewRouteUserController(UserController)

	BookController = controllers.NewBookController(bookService)
	BookRouteController = routes.NewBookRouteController(BookController)

	server = gin.New()         // Create a new Gin engine without default middleware
	server.Use(gin.Logger())   // Add Logger middleware explicitly
	server.Use(gin.Recovery()) // Add Recovery middleware explicitly
}

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer func() {
		if mongoclient != nil {
			mongoclient.Disconnect(ctx)
		}
	}()

	value := ""
	if redisclient != nil {
		value, err = redisclient.Get("key").Result()
		if err == redis.Nil {
			fmt.Println("key: test does not exist")
		} else if err != nil {
			panic(err)
		}
	} else {
		log.Println("Redis client is nil, skipping key retrieval")
	}

	server.Static("/static", "./static")

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8080", "http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/library")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	AuthRouteController.AuthRoute(router, userService)
	UserRouteController.UserRoute(router, userService)
	BookRouteController.BookRoute(router, bookService)

	log.Println("Registered routes:")
	for _, route := range server.Routes() {
		log.Printf("Method: %s, Path: %s", route.Method, route.Path)
	}

	log.Fatal(server.Run(":" + config.Port))
}
