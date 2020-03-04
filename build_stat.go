package main

import (
	"context"
	"fmt"
	"log"

	"time"

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
	// userRating := database.Collection("rating")
	// cur, err := usersCollection.Find(ctx, bson.M{})

	// defer cur.Close(ctx)
	// i := 0
	// countUsers, _ := usersCollection.CountDocuments(ctx, bson.M{})
	// var throttler = make(chan bool, MAX_CONCURENCY)
	// wg := new(sync.WaitGroup)
	// for cur.Next(ctx) {
	// 	var result models.User
	// 	i += 1
	// 	fmt.Printf("%v/%v\n", i, countUsers)
	// 	err := cur.Decode(&result)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	throttler <- true
	// 	wg.Add(1)

	// 	go func(result models.User) {
	// 		defer wg.Done()
	// 		number, err := userGamesCollection.CountDocuments(ctx, bson.M{"userid": result.ID})
	// 		options := options.FindOneAndReplace()
	// 		options.SetUpsert(true)
	// 		userRating.FindOneAndReplace(ctx, bson.M{"_id": result.ID}, bson.M{"_id": result.ID, "games_played": number}, options)
	// 		fmt.Println(number)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		<-throttler
	// 	}(result)
	// }
	// userRating.Indexes().CreateOne(ctx, mongo.IndexModel{
	// 	Keys: bson.M{
	// 		"game_played": 1,
	// 	}, Options: nil,
	// })

	res, err := userGamesCollection.Distinct(context.Background(), "gametype", bson.M{})
	fmt.Println(res)
	fmt.Println(err)
	for _, gametype := range res {
		fmt.Println(gametype)
		pipeline := []bson.D{
			bson.D{{"$match", bson.D{{"gametype", "11"}}}},
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
		cur, err := usersCollection.Aggregate(context.Background(), pipeline, opts)
		defer cur.Close(ctx)
		fmt.Println(cur)
		fmt.Println(err)

		var results []bson.M
		cur.All(ctx, &results)

		fmt.Println(results)

	}
}
