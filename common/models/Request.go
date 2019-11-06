package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// RequestModel - MongoDB mapping
type RequestModel struct {
	ID           bson.ObjectId          `bson:"_id" json:"_id"`
	ExperimentID bson.ObjectId          `bson:"experiment_id" json:"experiment_id"`
	RequestType  string                 `bson:"request_type" json:"request_type"`
	Status       string                 `bson:"status" json:"status"`
	Rules        map[string]interface{} `bson:"rules" json:"rules"`
	CreatedAt    time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time              `bson:"updated_at" json:"updated_at"`
}
