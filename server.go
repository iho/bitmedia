package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/iho/bitmedia/handlers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database("bitmedia")
	env := handlers.Env{Db: database}

	// TODO move to separate function
	usersCollection := database.Collection("users")
	userGamesCollection := database.Collection("user_games")
	buildIndixes(ctx, usersCollection, userGamesCollection)

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://postwoman.io"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/users/:id/stats", env.UserStats)

	router.GET("/users", env.ListUsers)

	router.GET("/games", env.ListGames)
	router.GET("/games/stats", env.GameStats)
	router.GET("/rating", env.UserRating)

	router.Run()
}

func buildIndixes(ctx context.Context, usersCollection *mongo.Collection, userGamesCollection *mongo.Collection) {
	_, err := usersCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.M{
			"birth_date": -1,
		},
	}, {
		Keys: bson.M{
			"gender": -1,
		}, Options: nil,
	}, {
		Keys: bson.M{
			"city": -1,
		}, Options: nil,
	}, {
		Keys: bson.M{
			"country": -1,
		}, Options: nil,
	}, {
		Keys: bson.M{
			"lastname": -1,
		}, Options: nil,
	}, {
		Keys: bson.M{
			"email": -1,
		}, Options: nil,
	}})
	if err != nil {
		log.Fatal(err)
	}

	_, err = userGamesCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.M{
			"created": -1,
		}, Options: nil,
	}, {
		Keys: bson.M{
			"userid": -1,
		}, Options: nil,
	}, {
		Keys: bson.M{
			"gametype": -1,
		}, Options: nil,
	}, {
		Keys: bson.M{
			"pointsgained": -1,
		}, Options: nil,
	}, {
		Keys: bson.M{
			"winstatus": -1,
		}, Options: nil,
	}})
	if err != nil {
		log.Fatal(err)
	}
}
