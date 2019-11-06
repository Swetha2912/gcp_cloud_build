package models

import (
	"time"
	// "go.mongodb.org/mongo-driver/mongo"
	// "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// DatasourceModel - MongoDB mapping
type DatasourceModel struct {
	ID             bson.ObjectId `bson:"_id" json:"_id"`
	UserID         bson.ObjectId `bson:"user_id,omitempty" json:"user_id"`
	DatasourceType string        `bson:"datasource_type" json:"datasource_type"` // folder | s3 | gdrive etc
	AccessKey      string        `bson:"access_key" json:"access_key"`
	AccessSecret   string        `bson:"access_secret" json:"access_secret"`
	BucketName     string        `bson:"bucket_name" json:"bucket_name"`
	FilterPath     string        `bson:"filter_path,omitempty" json:"filter_path,omitempty"`
	Region         string        `bson:"region" json:"region"`
	IsActive       bool          `bson:"is_active" json:"is_active"`
	CreatedAt      time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time     `bson:"updated_at" json:"updated_at"`
}
