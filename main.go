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
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

var recipesHandler *handlers.RecipesHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	recipesHandler = handlers.NewRecipesHandler(ctx, collection)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	recipesIdRoute := "recipes/:id"
	router.PUT(recipesIdRoute, recipesHandler.UpdateRecipeHandler)
	router.DELETE(recipesIdRoute, recipesHandler.DeleteRecipeHandler)
	router.GET(recipesIdRoute, recipesHandler.GetOneRecipeHandler)
	//router.GET("/recipes/search", SearchRecipesHandler)
	router.Run()
}
