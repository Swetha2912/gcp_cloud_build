package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// ExperimentImageModel -- mongodb mapping of Image inside Experiment
type ExperimentImageModel struct {
	ID             bson.ObjectId `bson:"_id" json:"_id"`
	ExperimentID   bson.ObjectId `bson:"experiment_id" json:"experiment_id,omitempty"`
	DatapointID    bson.ObjectId `bson:"datapoint_id" json:"datapoint_id,omitempty"`
	OriginalName   string        `bson:"original_name" json:"original_name,omitempty"`
	AnnotatorID    bson.ObjectId `bson:"annotator_id,omitempty" json:"annotator_id,omitempty"`
	Width          float32       `bson:"width" json:"width,omitempty"`
	Height         float32       `bson:"height" json:"height,omitempty"`
	DatasetID      bson.ObjectId `bson:"dataset_id" json:"dataset_id,omitempty"`
	DatasetImageID bson.ObjectId `bson:"dataset_image_id" json:"dataset_image_id,omitempty"`
	DatasourceType string        `bson:"datasource_type,omitempty" json:"datasource_type,omitempty"`
	ImageKey       string        `bson:"image_key" json:"image_key,omitempty"`
	IsAnnotated    bool          `bson:"is_annotated" json:"is_annotated,omitempty"`
	IsReviewed     bool          `bson:"is_reviewed" json:"is_reviewed,omitempty"`
	Status         string        `bson:"status" json:"status,omitempty"` // pending, annotate,
	CreatedAt      time.Time     `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt      time.Time     `bson:"updated_at" json:"updated_at,omitempty"`
	URL            string        `bson:"url,omitempty" json:"url,omitempty"`
}
