package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/iho/bitmedia/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GameStatsQueryParams struct {
	GameType  string
	Skip      int64
	Limit     int64
	DateStart time.Time `time_format:"2006-01-02"`
	DateEnd   time.Time `time_format:"2006-01-02"`
}

func (e *Env) GameStats(c *gin.Context) {

	ctx := c.Request.Context()
	var params GameStatsQueryParams
	err := c.ShouldBindWith(&params, binding.Query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error"})
	}
	var limit int64 = 100
	if params.Limit != 0 {
		limit = params.Limit
	}
	options := options.Find()
	options.SetLimit(limit)
	options.SetSkip(params.Skip)
	options.SetSort(bson.D{bson.E{Key: "date", Value: 1}, bson.E{Key: "gametype", Value: 1}})

	dbParams := bson.M{}

	if params.GameType != "" {
		dbParams["game_type"] = params.GameType
	}

	createdDate := bson.M{}
	if !params.DateEnd.IsZero() {
		createdDate["$lte"] = params.DateEnd
	}

	if !params.DateStart.IsZero() {
		createdDate["$gt"] = params.DateStart
	}
	if len(createdDate) != 0 {
		dbParams["date"] = createdDate
	}
	userStats := e.Db.Collection("stats")
	cur, err := userStats.Find(ctx, dbParams, options)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var results []models.GameStats
	err = cur.All(ctx, &results)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"resulsts": results})

}

type UserRatingQueryParams struct {
	Limit int64
	Skip  int64
}

func (e *Env) UserRating(c *gin.Context) {
	ctx := c.Request.Context()
	usersCollection := e.Db.Collection("users")
	// ctx := c.Request.Context()
	var params UserRatingQueryParams
	err := c.ShouldBindWith(&params, binding.Query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	var limit int64 = 100
	if params.Limit != 0 {
		limit = params.Limit
	}
	options := options.Find()
	options.SetLimit(limit)
	options.SetSkip(params.Skip)
	options.SetSort(bson.M{"games_played": -1})
	cur, err := usersCollection.Find(ctx, bson.M{}, options)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var results []models.User
	err = cur.All(ctx, &results)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rating": results})
}
