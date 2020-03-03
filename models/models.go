package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email         string             `json:"email" validate:"required,email"`
	LastName      string             `json:"last_name" validate:"required"`
	Country       string             `json:"country" validate:"required"`
	City          string             `json:"city" validate:"required"`
	Gender        string             `json:"gender" validate:"required"`
	BirthDate     string             `json:"birth_date,omitempty" bson:"-" validate:"required"`
	BirthDateTime time.Time          `json:"-" bson:"birth_date"`
}

type GameResult struct {
	PointsGained string             `json:"points_gained"`
	WinStatus    string             `json:"win_status"`
	GameType     string             `json:"game_type"`
	Created      string             `json:"created" bson:"-"`
	CreatedTime  time.Time          `json:"-" bson:"created"`
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id"`
}
