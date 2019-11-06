package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// BookmarkModel - mongodb mapping of BookmarkModel
type BookmarkModel struct {
	ID           bson.ObjectId `bson:"_id" json:"_id"`
	Title        string        `bson:"title" json:"title"`
	Description  string        `bson:"description" json:"description"`
	ExperimentID bson.ObjectId `bson:"experiment_id,omitempty" json:"experiment_id,omitempty"`
	ProjectID    bson.ObjectId `bson:"project_id,omitempty" json:"project_id,omitempty"`
	Source       string        `bson:"source" json:"source"` // filter | handpick | bulk
	BookmarkType string        `bson:"bookmark_type" json:"bookmark_type"`
	IsActive     bool          `bson:"is_active" json:"is_active"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at" json:"updated_at"`
}
