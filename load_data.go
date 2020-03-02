package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"log"
	"os"
	"time"

	"github.com/iho/bitmedia/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserJSON struct {
	Objects []models.User `json:"objects"`
}

type GameResultJSON struct {
	Objects []models.GameResult `json:"objects"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database("bitmedia")

	usersCollection := database.Collection("users")
	var usersJSON UserJSON
	jsonFile, err := os.Open("data/users_go.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &usersJSON)
	if err != nil {
		log.Fatal(err)
	}

	var usersInterfacesArray []interface{} = make([]interface{}, len(usersJSON.Objects))
	for i, user := range usersJSON.Objects {
		usersInterfacesArray[i] = user
	}

	_, err = usersCollection.InsertMany(ctx, usersInterfacesArray)

	if err != nil {
		log.Fatal(err)
	}

	userGamesCollection := database.Collection("user_games")
	var gamesResultJSON GameResultJSON
	jsonFile, err = os.Open("data/games.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	byteValue, _ = ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &gamesResultJSON)
	if err != nil {
		log.Fatal(err)
	}
	var userGamesInterfacesArray []interface{} = make([]interface{}, len(usersJSON.Objects))
	for i, user := range usersJSON.Objects {
		userGamesInterfacesArray[i] = user
	}
	_, err = userGamesCollection.InsertMany(ctx, userGamesInterfacesArray)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Files have been loaded successfully")
}
