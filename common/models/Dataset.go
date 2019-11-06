package models

import (
	"time"
	// "go.mongodb.org/mongo-driver/mongo"
	// "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// DatasetExportModel - MongoDB mapping
type DatasetExportModel struct {
	ID           bson.ObjectId `bson:"_id" json:"_id"`
	URL          string        `bson:"url" json:"url"`
	Classes      []string      `bson:"classes" json:"classes"`
	ImgCount     int           `bson:"img_count" json:"img_count"`
	AnnCount     int           `bson:"ann_count" json:"ann_count"`
	ExportFormat string        `bson:"export_format" json:"export_format"`
	TrainSplit   float64       `bson:"train_split,omitempty" json:"train_split,omitempty"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at" json:"updated_at"`
}

// DatasetModel - MongoDB mapping
type DatasetModel struct {
	ID                   bson.ObjectId        `bson:"_id" json:"_id"`
	UserID               bson.ObjectId        `bson:"user_id,omitempty" json:"user_id"`
	Title                string               `bson:"title" json:"title"`
	Description          string               `bson:"description" json:"description"`
	DatasetType          string               `bson:"dataset_type,omitempty" json:"dataset_type,omitempty"`
	DatasourceType       string               `bson:"datasource_type" json:"datasource_type"` // folder | s3 | gdrive etc
	DatasourceID         bson.ObjectId        `bson:"datasource_id,omitempty" json:"datasource_id,omitempty"`
	CredentialID         bson.ObjectId        `bson:"credential_id,omitempty" json:"credential_id,omitempty"`
	Region               string               `bson:"region,omitempty" json:"region,omitempty"`
	DatasourceBucketName string               `bson:"datasource_bucket_name,omitempty" json:"datasource_bucket_name,omitempty"`
	DatasourceFilterPath string               `bson:"datasource_filter_path,omitempty" json:"datasource_filter_path,omitempty"`
	LastSyncDate         time.Time            `bson:"last_sync_date,omitempty" json:"last_sync_date,omitempty"`
	Exports              []DatasetExportModel `bson:"exports" json:"exports"`
	IsActive             bool                 `bson:"is_active" json:"is_active"`
	CreatedAt            time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt            time.Time            `bson:"updated_at" json:"updated_at"`
}
