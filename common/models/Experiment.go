package models

import (
	"time"

	// "go.mongodb.org/mongo-driver/mongo"
	// "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// ExperimentExportModel - MongoDB mapping
type ExperimentExportModel struct {
	ID           bson.ObjectId `bson:"_id" json:"_id"`
	URL          string        `bson:"url" json:"url"`
	Classes      []string      `bson:"classes" json:"classes"`
	ImgCount     int           `bson:"img_count" json:"img_count"`
	AnnCount     int           `bson:"ann_count" json:"ann_count"`
	ExportFormat string        `bson:"export_format" json:"export_format"`
	TrainSplit   float64       `bson:"train_split,omitempty" json:"train_split,omitempty"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at" json:"updated_at"`
}

// InferenceRequest - MongoDB mapping
type InferenceRequest struct {
	ID              bson.ObjectId          `bson:"_id" json:"_id"`
	ModelID         bson.ObjectId          `bson:"model_id" json:"model_id"`
	ExperimentID    bson.ObjectId          `bson:"experiment_id" json:"experiment_id"`
	Thresholds      map[string]interface{} `bson:"thresholds" json:"thresholds"`
	Status          string                 `bson:"status" json:"status"`
	InferenceOutput interface{}            `bson:"inference_output,omitempty" json:"inference_output,omitempty"`
	CreatedAt       time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time              `bson:"updated_at" json:"updated_at"`
}

// ExperimentModel - MongoDB mapping
type ExperimentModel struct {
	ID           bson.ObjectId `bson:"_id" json:"_id"`
	ProjectID    bson.ObjectId `bson:"project_id" json:"project_id"`
	UserID       bson.ObjectId `bson:"user_id,omitempty" json:"user_id"`
	Title        string        `bson:"title" json:"title"`
	Architecture string        `bson:"architecture" json:"architecture"`
	Framework    string        `bson:"framework" json:"framework"`
	ProblemType  string        `bson:"problem_type" json:"problem_type"`
	Description  string        `bson:"description" json:"description"`
	IsActive     bool          `bson:"is_active" json:"is_active"`
	Status       string        `bson:"status" json:"status"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at" json:"updated_at"`
	// Classes      []string                `bson:"classes" json:"classes"`
	Datapoints []DatapointModel        `bson:"datapoints" json:"datapoints"`
	Exports    []ExperimentExportModel `bson:"exports" json:"exports"`
}
