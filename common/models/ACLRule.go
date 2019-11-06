package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// ACLRule -- This stores the userlevel or group-level or team-level -- read, write or delete access
type ACLRule struct {
	ID          bson.ObjectId `bson:"_id" json:"_id"`
	UserID      bson.ObjectId `bson:"user_id,omitempty" json:"user_id,omitempty"`
	TeamID      bson.ObjectId `bson:"team_id,omitempty" json:"team_id,omitempty"`
	AccountType string        `bson:"account_type,omitempty" json:"account_type,omitempty"`
	Permissions []string      `bson:"permissions" json:"permissions"` // can be read, edit, delete
	RuleType    string        `bson:"rule_type" json:"rule_type"`     // module or resource -- can manager access this resource (or)
	IsActive    bool          `bson:"is_active" json:"is_active"`     // can annotator access AnnotationUI

	// for resource type
	ResourceID   bson.ObjectId `bson:"resource_id,omitempty" json:"resource_id"`
	ResourceName string        `bson:"resource_name" json:"resource_name"` // dataset, experiment, project
	ResourceType string        `bson:"resource_type" json:"resource_type"` // module or resource

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
