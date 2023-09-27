// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/IjatAyam/recipes-api.
//
//	Schemes: http
//	Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//	Contact: Wan Muhammad Izzat <ijat191999@gmail.com> https://jatyam.com
//
//	Consumes:
//		- application/json
//
//	Produces:
//		- application/json
//
// swagger:meta
package main

import (
	"context"
	"github.com/IjatAyam/recipes-api/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

var authHandler *handlers.AuthHandler
var recipesHandler *handlers.RecipesHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	log.Println(status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)

	collectionUsers := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
}

func SetupServer() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/recipes", recipesHandler.ListRecipesHandler)

	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)

	authorized := router.Group("/recipes")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.POST("", recipesHandler.NewRecipeHandler)
		recipesIdRoute := "/recipes/:id"
		authorized.PUT(recipesIdRoute, recipesHandler.UpdateRecipeHandler)
		authorized.DELETE(recipesIdRoute, recipesHandler.DeleteRecipeHandler)
		authorized.GET(recipesIdRoute, recipesHandler.GetOneRecipeHandler)
	}

	return router
}

func main() {
	//SetupServer().RunTLS(":443", "certs/localhost.crt", "certs/localhost.key")
	SetupServer().Run()
}
