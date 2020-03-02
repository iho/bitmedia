package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"

	"log"
	"os"
	"time"

	"github.com/iho/bitmedia/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserJSON struct {
	Objects []models.User `json:"objects"`
}

type GameResultJSON struct {
	Objects []models.GameResult `json:"objects"`
}

const (
	MINIMUM_SIZE                  int = 5000
	MAXIMUM_SIZE                  int = MINIMUM_SIZE * 2
	NUMBER_OF_PARALLEL_INSERTIONS int = 4
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
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

	insertionResult, err := usersCollection.InsertMany(ctx, usersInterfacesArray)

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
	rand.Seed(time.Now().Unix())
	ch := make(chan bool, NUMBER_OF_PARALLEL_INSERTIONS)
	for _, userID := range insertionResult.InsertedIDs {
		primitiveUserID := userID.(primitive.ObjectID)
		go InsertGames(ctx, userGamesCollection, primitiveUserID, &gamesResultJSON, ch)
		ch <- true
	}

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Files have been loaded successfully")
}

func InsertGames(ctx context.Context, collection *mongo.Collection, UserID primitive.ObjectID, games *GameResultJSON, ch chan bool) {
	<-ch
	quantity := rand.Intn(MAXIMUM_SIZE-MINIMUM_SIZE) + MINIMUM_SIZE
	var userGamesInterfacesArray []interface{} = make([]interface{}, quantity)
	for i := 0; i < quantity; i++ {
		game := games.Objects[rand.Intn(len(games.Objects))]
		game.UserID = UserID
		userGamesInterfacesArray[i] = game
	}
	_, err := collection.InsertMany(ctx, userGamesInterfacesArray)
	if err != nil {
		log.Fatal(err)
	}
}
