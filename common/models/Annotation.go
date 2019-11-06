package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// AnnotationModel - mongodb mapping of all types annotatations
type AnnotationModel struct {
	ID                 bson.ObjectId `bson:"_id" json:"_id"`
	ExperimentID       bson.ObjectId `bson:"experiment_id" json:"experiment_id"`
	ImageID            bson.ObjectId `bson:"image_id" json:"image_id"`
	Source             string        `bson:"source" json:"source"` // human | machine | review |
	Type               string        `bson:"type" json:"type"`     //segmentation or detection
	Class              string        `bson:"cls" json:"cls"`
	Coord              interface{}   `bson:"coords" json:"coords"`
	IsActive           bool          `bson:"is_active" json:"is_active"`
	InferenceRequestID bson.ObjectId `bson:"inference_request_id,omitempty" json:"inference_request_id,omitempty"`
	Confidence         float64       `bson:"confidence" json:"confidence"` // active learning confidence is stored here
	CreatedAt          time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time     `bson:"updated_at" json:"updated_at"`
}
