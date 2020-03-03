package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"

	"os"
	"sync"
	"time"

	"github.com/iho/bitmedia/models"

	"github.com/jeffjen/datefmt"

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
	MINIMUM_SIZE   int = 5000
	MAXIMUM_SIZE   int = MINIMUM_SIZE * 2
	MAX_CONCURENCY int = 16
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
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

	const dateFormat = "%A, %B %-d, %Y %I:%M %p"

	for i, user := range usersJSON.Objects {
		t, err := datefmt.Strptime(dateFormat, user.BirthDate)
		if err != nil {
			fmt.Println(err)
			t = time.Now()
		}

		user.BirthDateTime = t
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
	parseGameDates(&gamesResultJSON)

	rand.Seed(time.Now().Unix())
	wg := new(sync.WaitGroup)
	wg.Add(len(insertionResult.InsertedIDs))
	var throttler = make(chan bool, MAX_CONCURENCY)
	for _, userID := range insertionResult.InsertedIDs {
		primitiveUserID := userID.(primitive.ObjectID)
		throttler <- true
		go func() {
			defer wg.Done()
			InsertGames(ctx, userGamesCollection, primitiveUserID, &gamesResultJSON)
			<-throttler
		}()

	}

	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
	fmt.Println("Files have been loaded successfully")
}

func InsertGames(ctx context.Context, collection *mongo.Collection, UserID primitive.ObjectID, games *GameResultJSON) {
	quantity := rand.Intn(MAXIMUM_SIZE-MINIMUM_SIZE) + MINIMUM_SIZE
	var userGamesInterfacesArray []interface{} = make([]interface{}, quantity)
	for i := 0; i < quantity; i++ {
		game := games.Objects[rand.Intn(len(games.Objects))]
		game.UserID = UserID
		userGamesInterfacesArray[i] = game
	}
	res, err := collection.InsertMany(ctx, userGamesInterfacesArray)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(res.InsertedIDs))
}

func parseGameDates(games *GameResultJSON) {
	for i := 0; i < len(games.Objects); i++ {
		game := &games.Objects[i]
		t, err := time.Parse("1/2/2006 15:04 PM", game.Created)
		if err != nil {
			fmt.Println(err)
			t = time.Now()
		}
		game.CreatedTime = t
		games.Objects[i] = *game
	}
}
