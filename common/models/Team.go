package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// TeamUser - MongoDB mapping
type TeamUser struct {
	ID        bson.ObjectId `bson:"_id" json:"_id"`
	TeamID    bson.ObjectId `bson:"team_id" json:"team_id"`
	UserID    bson.ObjectId `bson:"user_id" json:"user_id"`
	Role      string        `bson:"role" json:"role"`
	IsActive  bool          `bson:"is_active" json:"is_active"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

// TeamModel - MongoDB mapping
type TeamModel struct {
	ID          bson.ObjectId `bson:"_id" json:"_id"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
	IsActive    bool          `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
}
