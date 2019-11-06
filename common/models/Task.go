package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// TaskModel - MongoDB mapping
type TaskModel struct {
	ID           bson.ObjectId          `bson:"_id" json:"_id"`
	ExperimentID bson.ObjectId          `bson:"experiment_id" json:"experiment_id"`
	RequestID    bson.ObjectId          `bson:"request_id" json:"request_id"`
	AnnotatorID  bson.ObjectId          `bson:"annotator_id" json:"annatator_id"`
	ProblemType  string                 `bson:"problem_type" json:"problem_type"`
	Status       string                 `bson:"status" json:"status"`
	Rules        map[string]interface{} `bson:"rules" json:"rules"`
	CreatedAt    time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time              `bson:"updated_at" json:"updated_at"`
}
