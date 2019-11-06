package models

import (
	"time"
	// "go.mongodb.org/mongo-driver/mongo"
	// "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// TmpModel - MongoDB mapping -- added to inference request during traing process
type TmpModel struct {
	ID         bson.ObjectId `bson:"_id" json:"_id"`
	SessionID  string        `json:"session_id"`
	ModelPath  string        `json:"string"`
	Accuracy   float64       `bson:"accuracy" json:"accuracy"`
	Loss       float64       `bson:"loss" json:"loss"`
	Steps      int           `bson:"steps" json:"steps"`
	IsUpgraded bool          `bson:"is_upgraded" json:"is_upgraded"`
	IsActive   bool          `bson:"is_active" json:"is_active"`
	CreatedAt  time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time     `bson:"updated_at" json:"updated_at"`
}

// TrainedModel - MongoDB mapping
type TrainedModel struct {
	ID           bson.ObjectId       `bson:"_id" json:"_id"`
	Title        string              `bson:"title" json:"title"`
	ProblemType  string              `bson:"problem_type" json:"problem_type"`
	ExperimentID bson.ObjectId       `bson:"experiment_id,omitempty" json:"experiment_id,omitempty"`
	ProjectID    bson.ObjectId       `bson:"project_id,omitempty" json:"project_id"`
	UserID       bson.ObjectId       `bson:"user_id,omitempty" json:"user_id"`
	Architecture string              `bson:"architecture" json:"architecture"`
	Accuracy     float64             `bson:"accuracy" json:"accuracy"`
	Variant      string              `bson:"variant" json:"variant"`
	Loss         float64             `bson:"loss" json:"loss"`
	Steps        int                 `bson:"steps" json:"steps"`
	IsActive     bool                `bson:"is_active" json:"is_active"`
	IsFav        bool                `bson:"is_fav" json:"is_fav"`
	Source       string              `bson:"source" json:"source"`
	Classes      []string            `bson:"classes" json:"classes"`
	Files        []map[string]string `bson:"files" json:"files"`
	CreatedAt    time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time           `bson:"updated_at" json:"updated_at"`
	Status       string              `bson:"status" json:"status"`
}
