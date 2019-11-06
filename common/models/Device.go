package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// DeviceModel - MongoDB mapping
type DeviceModel struct {
	ID           bson.ObjectId `bson:"_id" json:"_id"`
	UserID       bson.ObjectId `bson:"user_id,omitempty" json:"user_id"`
	Name         string        `bson:"name" json:"name"`
	Type         string        `bson:"type" json:"type"` //gpu
	InstanceID   string        `bson:"instance_id" json:"instance_id"`
	CredentialID bson.ObjectId `bson:"credential_id,omitempty" json:"credential_id,omitempty"`
	Status       string        `bson:"status" json:"status"`
	IsActive     bool          `bson:"is_active" json:"is_active"`
	State        bool          `bson:"state" json:"state"`
	PublicIP     string        `bson:"public_ip" json:"public_ip"`
	LastSyncDate time.Time     `bson:"last_sync_date,omitempty" json:"last_sync_date,omitempty"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at" json:"updated_at"`
}
