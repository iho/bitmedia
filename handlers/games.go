package handlers

import (
	"fmt"
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
	After        time.Time `time_format:"2006-01-02T15:04:05"`
}

func (e *Env) ListGames(c *gin.Context) {
	var params ListsGamesQueryParams
	var gameResults []models.GameResult
	var limit int64 = 100
	if err := c.ShouldBindWith(&params, binding.Query); err == nil {
		userID, err := primitive.ObjectIDFromHex(params.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx := c.Request.Context()
		usersCollection := e.Db.Collection("user_games")
		createdDate := bson.M{}
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
		if params.UserID != "" {
			dbParams["userID"] = userID
		}

		if (params.CreatedStart != time.Time{}) {
			if (params.After == time.Time{}) {
				createdDate["$gt"] = params.CreatedStart
			} else {
				bigger := params.CreatedStart
				if params.After.After(params.CreatedStart) {
					bigger = params.After
				}

				createdDate["$gt"] = bigger

			}
		} else if (params.After != time.Time{}) {
			createdDate["$gt"] = params.After
		}
		if (params.CreatedEnd != time.Time{}) {
			createdDate["$lte"] = params.CreatedEnd

		}
		fmt.Println(createdDate)

		if createdDate != nil {
			dbParams["created"] = createdDate
		}
		fmt.Println(dbParams)

		findOptions := options.Find()
		findOptions.SetLimit(limit)
		findOptions.SetSort(bson.M{"created": -1})
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
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
