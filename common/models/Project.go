package models

import (
	"time"
	// "go.mongodb.org/mongo-driver/mongo"
	// "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// ClassModel --
type ClassModel struct {
	ID           bson.ObjectId `bson:"_id" json:"_id"`
	Name         string        `bson:"name" json:"name"`
	ExperimentID bson.ObjectId `bson:"experiment_id" json:"experiment_id"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at" json:"updated_at"`
}

// ProjectModel -- mongodb mapping
type ProjectModel struct {
	ID          bson.ObjectId `bson:"_id" json:"_id"`
	UserID      bson.ObjectId `bson:"user_id,omitempty" json:"user_id"`
	Title       string        `bson:"title" json:"title"`
	TeamID      bson.ObjectId `bson:"team_id,omitempty" json:"team_id,omitempty"`
	Description string        `bson:"description" json:"description"`
	IsActive    bool          `bson:"is_active" json:"is_active"`
	Status      string        `bson:"status" json:"status"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
	Classes     []string      `bson:"classes" json:"classes"`
}
