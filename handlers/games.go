package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/iho/bitmedia/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListsGamesQueryParams struct {
	PointsGained string
	WinStatus    string
	GameType     string
	UserID       string
	CreatedStart time.Time `time_format:"2006-01-02T15:04:05"`
	CreatedEnd   time.Time `time_format:"2006-01-02T15:04:05"`
	After        string
}

func (e *Env) ListGames(c *gin.Context) {
	var params ListsGamesQueryParams
	var gameResults []models.GameResult
	var limit int64 = 100
	err := c.ShouldBindWith(&params, binding.Query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	usersCollection := e.Db.Collection("user_games")
	dbParams := bson.M{}
	if params.PointsGained != "" {
		dbParams["pointsgained"] = params.PointsGained

	}
	if params.WinStatus != "" {
		dbParams["winstatus"] = params.WinStatus
	}
	if params.GameType != "" {
		dbParams["gametype"] = params.GameType
	}

	createdDate := bson.M{}
	if !params.CreatedEnd.IsZero() {
		createdDate["$lte"] = params.CreatedEnd
	}
	if !params.CreatedStart.IsZero() {
		createdDate["$get"] = params.CreatedStart
	}
	if len(createdDate) != 0 {
		dbParams["created"] = createdDate
	}

	if params.After != "" {
		gameID, err := primitive.ObjectIDFromHex(params.After)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		dbParams["_id"] = bson.M{"$gt": gameID}
	}
	if params.UserID != "" {
		userID, err := primitive.ObjectIDFromHex(params.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		dbParams["userID"] = userID
	}


	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.M{"_id": 1})
	cur, err := usersCollection.Find(ctx, dbParams, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result models.GameResult
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		// TODO move in functions
		result.Created = result.CreatedTime.Format("2006-01-02T15:04:05")
		result.CreatedTime = time.Time{}
		gameResults = append(gameResults, result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"gameResults": gameResults})

}
