package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// BookmarkImageModel -- mongodb mapping of Image inside Bookmark
type BookmarkImageModel struct {
	ID           bson.ObjectId `bson:"_id" json:"_id"`
	BookmarkID   bson.ObjectId `bson:"bookmark_id" json:"bookmark_id"`
	DatasetID    bson.ObjectId `bson:"dataset_id" json:"dataset_id"`
	ExperimentID bson.ObjectId `bson:"experiment_id" json:"experiment_id"`
	Width        float32       `bson:"width" json:"width"`
	Height       float32       `bson:"height" json:"height"`
	ImageKey     string        `bson:"image_key" json:"image_key"`
	UIHash       string        `bson:"ui_hash" json:"ui_hash"`
	BackendHash  string        `bson:"backend_hash" json:"backend_hash"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at" json:"updated_at"`
}
