package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

//Thumbnail -- thumbnails structure
type Thumbnail struct {
	Size        string `bson:"size"`
	IsGenerated bool   `bson:"is_generated"`
}

// DatasetImageModel -- mongodb mapping of Image inside Experiment
type DatasetImageModel struct {
	ID           bson.ObjectId `bson:"_id" json:"_id"`
	DatasetID    bson.ObjectId `bson:"dataset_id" json:"dataset_id"`
	OriginalName string        `bson:"original_name" json:"original_name"`
	Width        float32       `bson:"width" json:"width"`
	Height       float32       `bson:"height" json:"height"`
	ImageKey     string        `bson:"image_key" json:"image_key"`
	UIHash       string        `bson:"ui_hash" json:"ui_hash"`
	BackendHash  string        `bson:"backend_hash" json:"backend_hash"`
	// Thumbnails   []Thumbnail   `bson:"thumbnails"`
	DatasourceType     string    `bson:"datasource_type" json:"datasource_type"` // folder | s3 | gdrive etc
	ThumbnailProcessed bool      `bson:"thumbnail_processed" json:"thumbnail_processed"`
	ImageCaptured      time.Time `bson:"image_captured,omitempty" json:"image_captured,omitempty"`
	CreatedAt          time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time `bson:"updated_at" json:"updated_at"`
}
