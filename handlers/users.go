package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/iho/bitmedia/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (e *Env) UserStats(c *gin.Context) {
	id := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	usersCollection := e.Db.Collection("user_games")
	ctx := c.Request.Context()
	number, err := usersCollection.CountDocuments(ctx, bson.M{"userid": userID})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"gamesPlayed": number})

}

type ListsUsersQueryParams struct {
	Email          string
	LastName       string
	Country        string
	City           string
	Gender         string
	BirthDateStart time.Time `time_format:"2006-01-02"`
	BirthDateEnd   time.Time `time_format:"2006-01-02"`
	After          time.Time `time_format:"2006-01-02"`
}

type Env struct {
	Db *mongo.Database
}

func (e *Env) ListUsers(c *gin.Context) {
	var params ListsUsersQueryParams
	var users []models.User
	var limit int64 = 100
	err := c.ShouldBindWith(&params, binding.Query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	usersCollection := e.Db.Collection("users")

	dbParams := buildDBParams(params)
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.M{"birth_date": 1})
	cur, err := usersCollection.Find(ctx, dbParams, findOptions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result models.User
		err := cur.Decode(&result)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		// TODO move in functions
		result.BirthDate = result.BirthDateTime.Format("2006-01-02")
		result.BirthDateTime = time.Time{}
		users = append(users, result)
	}
	if err := cur.Err(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"users": users})

}

func buildDBParams(params ListsUsersQueryParams) bson.M {
	dbParams := bson.M{}
	if params.Email != "" {
		dbParams["email"] = params.Email
	}
	if params.LastName != "" {
		dbParams["lastname"] = params.LastName
	}
	if params.Country != "" {
		dbParams["country"] = params.Country
	}
	if params.City != "" {
		dbParams["city"] = params.City
	}
	if params.Gender != "" {
		dbParams["gender"] = params.Gender
	}
	birthDate := bson.M{}

	if !params.BirthDateEnd.IsZero() {
		birthDate["$lte"] = params.BirthDateEnd
	}

	bigger := params.BirthDateStart
	if params.After.After(params.BirthDateStart) {
		bigger = params.After
	}
	birthDate["$gt"] = bigger

	if len(birthDate) != 0 {
		dbParams["birth_date"] = birthDate
	}
	return dbParams
}
