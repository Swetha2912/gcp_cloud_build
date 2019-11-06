package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// CredentialModel - MongoDB mapping
type CredentialModel struct {
	ID           bson.ObjectId `bson:"_id" json:"_id"`
	UserID       bson.ObjectId `bson:"user_id,omitempty" json:"user_id"`
	Name         string        `bson:"name" json:"name"`
	AccessKey    string        `bson:"access_key" json:"access_key"`
	AccessSecret string        `bson:"access_secret" json:"access_secret"`
	Region       string        `bson:"region" json:"region"`
	IsActive     bool          `bson:"is_active" json:"is_active"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at" json:"updated_at"`
}
