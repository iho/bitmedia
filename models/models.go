package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Email         string    `json:"email"`
	LastName      string    `json:"last_name"`
	Country       string    `json:"country"`
	City          string    `json:"city"`
	Gender        string    `json:"gender"`
	BirthDate     string    `json:"birth_date" bson:"-"`
	BirthDateTime time.Time `json:"birth_date_time,omitempty" bson:"birth_date"`
}

type GameResult struct {
	PointsGained string             `json:"points_gained"`
	WinStatus    string             `json:"win_status"`
	GameType     string             `json:"game_type"`
	Created      string             `json:"created" bson:"-"`
	CreatedTime  time.Time          `json:"created_time,omitempty" bson:"created"`
	UserID       primitive.ObjectID `json:"user_id"`
}
