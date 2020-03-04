package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"time"

	"github.com/iho/bitmedia/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

const MAX_CONCURENCY int = 16

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database("bitmedia")

	usersCollection := database.Collection("users")
	userGamesCollection := database.Collection("user_games")

	cur, err := usersCollection.Find(ctx, bson.M{})

	defer cur.Close(ctx)
	i := 0
	countUsers, _ := usersCollection.CountDocuments(ctx, bson.M{})
	var throttler = make(chan bool, MAX_CONCURENCY)
	wg := new(sync.WaitGroup)

	options := options.FindOneAndUpdate()
	options.SetUpsert(true)

	for cur.Next(ctx) {
		var result models.User
		i += 1
		fmt.Printf("%v/%v\n", i, countUsers)
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		throttler <- true
		wg.Add(1)

		go func(result models.User) {
			defer wg.Done()
			number, err := userGamesCollection.CountDocuments(ctx, bson.M{"userid": result.ID})

			usersCollection.FindOneAndUpdate(ctx, bson.M{"_id": result.ID}, bson.M{"$set": bson.M{"games_played": number}}, options)
			fmt.Println(number)
			if err != nil {
				log.Fatal(err)
			}
			<-throttler
		}(result)
	}
	usersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"games_played": 1,
		}, Options: nil,
	})

}
