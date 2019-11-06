package models

import (
	"encoding/xml"
	"time"

	"github.com/globalsign/mgo/bson"
)

//XmlAnn -- xml structure
type XmlAnn struct {
	XMLName   xml.Name    `xml:"annotation"`
	XmlObject []XmlObject `xml:"object"`
	Filename  string      `xml:"filename"`
}

//XmlObject --  xml structure
type XmlObject struct {
	XMLName   xml.Name  `xml:"object"`
	Name      string    `xml:"name"`
	XmlBndbox XmlBndbox `xml:"bndbox"`
}

//XmlBndbox -- xml structure
type XmlBndbox struct {
	XMLName xml.Name `xml:"bndbox"`
	Xmin    string   `xml:"xmin"`
	Ymin    string   `xml:"ymin"`
	Xmax    string   `xml:"xmax"`
	Ymax    string   `xml:"ymax"`
}

// DatasetAnnotationModel - MongoDB mapping
type DatasetAnnotationModel struct {
	ID             bson.ObjectId `bson:"_id" json:"_id"`
	OriginalName   string        `bson:"original_name" json:"original_name"`
	DatasetID      bson.ObjectId `bson:"dataset_id" json:"dataset_id,omitempty"`
	DatasetImageID bson.ObjectId `bson:"dataset_image_id" json:"dataset_image_id,omitempty"`
	Class          string        `bson:"cls" json:"cls"`
	Coord          interface{}   `bson:"coords" json:"coords"`
	IsActive       bool          `bson:"is_active" json:"is_active"`
	ProblemType    string        `bson:"problem_type" json:"problem_type"`
	CreatedAt      time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time     `bson:"updated_at" json:"updated_at"`
}
