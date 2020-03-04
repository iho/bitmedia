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

type AggregateResponseID struct {
	Date     string `bson:"date"`
	Gametype string `bson:"gametype"`
}
type AggregateResponse struct {
	ID    AggregateResponseID `bson:"_id"`
	Count int                 `bson:"count"`
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database("bitmedia")

	userGamesCollection := database.Collection("user_games")
	userStats := database.Collection("stats")
	pipeline := []bson.D{
		bson.D{{"$limit", 100000000}},
		bson.D{{
			"$group", bson.M{
				"_id": bson.M{
					"date": bson.M{
						"$dateToString": bson.M{
							"format": "%Y-%m-%d",
							"date":   "$created",
						},
					},
					"gametype": "$gametype",
				},
				"count": bson.M{"$sum": 1},
			}}},
		bson.D{{"$sort", bson.M{"_id": 1}}},
	}
	opts := options.Aggregate().SetMaxTime(5 * time.Minute)
	cur, err := userGamesCollection.Aggregate(context.Background(), pipeline, opts)
	defer cur.Close(ctx)

	var throttler = make(chan bool, MAX_CONCURENCY)
	wg := new(sync.WaitGroup)
	for cur.Next(ctx) {
		var result AggregateResponse

		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
		throttler <- true
		wg.Add(1)

		go func(result AggregateResponse) {
			defer wg.Done()
			stat := models.GameStats{}
			stat.GameType = result.ID.Gametype
			stat.Count = result.Count

			t, err := time.Parse("2006-01-02", result.ID.Date)
			if err != nil {
				log.Fatal(err)

			}
			stat.Date = t
			options := options.FindOneAndUpdate()
			options.SetUpsert(true)
			userStats.FindOneAndUpdate(ctx, bson.M{"date": stat.Date, "gametype": stat.GameType}, bson.M{"$set": stat}, options)

			if err != nil {
				log.Fatal(err)
			}
			<-throttler
		}(result)
	}
	userStats.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"gametype": 1,
		}, Options: nil,
	})
	userStats.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"date": 1,
		}, Options: nil,
	})

}
