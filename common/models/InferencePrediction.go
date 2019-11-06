package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// InferencePredictionModel - mongodb mapping for inference requests predictions
type InferencePredictionModel struct {
	ID          bson.ObjectId  `bson:"_id" json:"_id"`
	InferenceID bson.ObjectId  `bson:"inference_request_id" json:"inference_request_id"`
	Class       string         `bson:"cls" json:"cls"`
	Coord       DetectionCoord `bson:"coords" json:"coords"`
	CreatedAt   time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `bson:"updated_at" json:"updated_at"`
}
