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
	c.JSON(http.StatusOK, gin.H{"number": number})

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
	if err := c.ShouldBindWith(&params, binding.Query); err == nil {
		ctx := c.Request.Context()
		usersCollection := e.Db.Collection("users")

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
		if (params.BirthDateStart != time.Time{} || params.BirthDateEnd != time.Time{} || params.After == time.Time{}) {
			birthDate := bson.M{}
			if (params.After != time.Time{}) {
				birthDate["$gt"] = params.After
				dbParams["birth_date"] = birthDate
				fmt.Println(params.After)
			}
			if (params.BirthDateEnd != time.Time{}) {
				birthDate["$lte"] = params.BirthDateEnd
				dbParams["birth_date"] = birthDate
				fmt.Println(params.BirthDateEnd)
			}
			if (params.BirthDateStart != time.Time{}) {
				if (params.After == time.Time{}) {
					birthDate["$gt"] = params.BirthDateStart
					dbParams["birth_date"] = birthDate
				} else {
					bigger := params.BirthDateStart
					if params.After.After(params.BirthDateStart) {
						bigger = params.After
					}
					birthDate["$gt"] = bigger
					dbParams["birth_date"] = birthDate
				}
			}

			fmt.Println(birthDate)

		}

		fmt.Println(dbParams)

		findOptions := options.Find()
		findOptions.SetLimit(limit)
		findOptions.SetSort(bson.M{"birth_date": -1})
		cur, err := usersCollection.Find(ctx, dbParams, findOptions)
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			var result models.User
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			// TODO move in functions
			result.BirthDate = result.BirthDateTime.Format("2006-01-02")
			result.BirthDateTime = time.Time{}
			users = append(users, result)
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"users": users})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// func UserStructLevelValidation(sl validator.StructLevel) {

// 	user := sl.Current().Interface().(User)

// 	if len(user.FirstName) == 0 && len(user.LastName) == 0 {
// 		sl.ReportError(user.FirstName, "fname", "FirstName", "fnameorlname", "")
// 		sl.ReportError(user.LastName, "lname", "LastName", "fnameorlname", "")
// 	}

// 	// plus can do more, even with different tag than "fnameorlname"
// }
