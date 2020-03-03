package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
)

type GameStatsQueryParams struct {
	GameType string
	Skip     int
	Limit    int
}

func (e *Env) GameStats(c *gin.Context) {
	usersCollection := e.Db.Collection("user_games")
	ctx := c.Request.Context()
	var params GameStatsQueryParams
	err := c.ShouldBindWith(&params, binding.Query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error"})
	}
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"gametype": params.GameType}},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"date": bson.M{
						"$dateToString": bson.M{
							"format": "%Y-%m-%d",
							"date": bson.M{
								"$toDate": "$created",
							},
						},
					},
					"gametype": "$gametype",
				},
				"count": bson.M{"$sum": 1},
			},
		},
		bson.M{"$sort": bson.M{"_id": 1}},
	}

	res, err := usersCollection.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"res": res})

}

type UserRatingQueryParams struct {
	Limit int
	Skip  int
}

func (e *Env) UserRating(c *gin.Context) {

	usersCollection := e.Db.Collection("user_games")
	ctx := c.Request.Context()
	var params UserRatingQueryParams
	err := c.ShouldBindWith(&params, binding.Query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	pipeline := []bson.M{
		bson.M{
			"$group": bson.M{
				"_id": "$userid",
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
		bson.M{
			"$sort": bson.M{
				"count": -1,
			},
		},
		bson.M{"$limit": params.Limit},
		bson.M{"$skip": params.Skip},
	}

	res, err := usersCollection.Aggregate(ctx, pipeline)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"res": res})
}
