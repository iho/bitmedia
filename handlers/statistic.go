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
	GameType     string
	Skip         int
	Limit        int
	CreatedStart time.Time `time_format:"2006-01-02"`
	CreatedEnd   time.Time `time_format:"2006-01-02"`
}

func (e *Env) GameStats(c *gin.Context) {
	// usersCollection := e.Db.Collection("user_games")
	// ctx := c.Request.Context()
	var params GameStatsQueryParams
	err := c.ShouldBindWith(&params, binding.Query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error"})
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"res": err.Error()})

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
