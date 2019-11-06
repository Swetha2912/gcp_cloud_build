package models

import (
	// "go.mongodb.org/mongo-driver/mongo"
	// "github.com/globalsign/mgo"
	"time"

	"github.com/globalsign/mgo/bson"
)

// DatapointModel -
type DatapointModel struct {
	ID            bson.ObjectId `bson:"_id" json:"_id"`
	LinkTo        string        `bson:"link_to" json:"link_to"` // class or empty=experiment level
	ExperimentID  bson.ObjectId `bson:"experiment_id,omitempty" json:"experiment_id,omitempty"`
	DatasetID     bson.ObjectId `bson:"dataset_id,omitempty" json:"dataset_id,omitempty"`
	BookmarkID    bson.ObjectId `bson:"bookmark_id,omitempty" json:"bookmark_id,omitempty"`
	ImagesFilter  []string      `bson:"images_filter,omitempty" json:"images_filter,omitempty"`
	ClassesFilter []string      `bson:"classes_filter,omitempty" json:"classes_filter,omitempty"`
	SourceFilter  []string      `bson:"source_filter,omitempty" json:"source_filter,omitempty"`
	Limit         int           `bson:"limit,omitempty" json:"limit,omitempty"`
	IsProcessed   bool          `bson:"is_processed,omitempty" json:"is_processed,omitempty"`
	RetainCoord   bool          `bson:"retain_coord,omitempty" json:"retain_coord,omitempty"`
	CreatedAt     time.Time     `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt     time.Time     `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
