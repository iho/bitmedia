package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
)

func (e *Env) GameStats(c *gin.Context) {
	// id := c.Param("gameid")
	// gameID, err := primitive.ObjectIDFromHex(id)

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	usersCollection := e.Db.Collection("user_games")
	ctx := c.Request.Context()

	pipeline := []bson.M{
		// bson.M{"$match": bson.M{"game_type": id},
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
	limit int
	skip  int
}

func (e *Env) UserRating(c *gin.Context) {

	usersCollection := e.Db.Collection("user_games")
	ctx := c.Request.Context()
	var params UserRatingQueryParams
	if err := c.ShouldBindWith(&params, binding.Query); err == nil {
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
			bson.M{"$limit": 100},
			bson.M{"$skip": 100},
		}

		res, err := usersCollection.Aggregate(ctx, pipeline)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"res": res})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}

// [{$match: {
//     gametype: "17"
//   }}, {$group: {
//     _id: {"$dateToString": {
//                   "format": "%Y-%m-%d",
//                   "date": {
//                       "$toDate": "$created"
//                   }
//               }},
//       "count": { "$sum": 1 }
//   }}, {$sort: {
//     _id: 1
//   }}]

//   [{$match: {
//     gametype: "17"
//   }}, {$group: {
//     _id: {"date": {"$dateToString": {
//                   "format": "%Y-%m-%d",
//                   "date": {
//                       "$toDate": "$created"
//                   }
//               }}, gametype: "$gametype"},
//       "count": { "$sum": 1 }
//   }}, {$sort: {
//     _id: 1
//   }}]
