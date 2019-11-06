package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// ExploreModel - MongoDB mapping
type ExploreModel struct {
	ID          bson.ObjectId       `bson:"_id" json:"_id"`
	Name        string              `bson:"name" json:"name"`
	Framework   string              `bson:"framework" json:"framework"`
	ModelPath   string              `bson:"model_path" json:"model_path"`
	ProjectType string              `bson:"project_type" json:"project_type"`
	Images      []string `bson:"images" json:"images"`
	IsActive    bool                `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time           `bson:"updated_at" json:"updated_at"`
}
