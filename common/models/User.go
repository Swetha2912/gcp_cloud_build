package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// UserModel ... mongodb mapping of all types annotatations
type UserModel struct {
	ID             bson.ObjectId `bson:"_id" json:"_id"`
	Name           string        `bson:"name" json:"name"`
	Email          string        `bson:"email" json:"email"` // human | machine | review |
	Phone          string        `bson:"phone" json:"phone"` //segmentation or detection
	Password       string        `bson:"password" json:"-"`
	ProfilePicture string        `bson:"profile_picture" json:"profile_picture"`
	Role           string        `bson:"role" json:"role"`
	IsActive       bool          `bson:"is_active" json:"is_active"`
	CreatedAt      time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time     `bson:"updated_at" json:"updated_at"`
}
