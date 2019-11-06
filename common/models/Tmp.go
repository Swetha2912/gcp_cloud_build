package models

import "github.com/globalsign/mgo/bson"

// ImageExport -- used during export of images
type ImageExport struct {
	Height       float32 `bson:"height" json:"height"`
	Width        float32 `bson:"width" json:"width"`
	ImageKey     string  `bson:"image_key" json:"image_key"`
	Originalname string  `bson:"original_name" json:"original_name"`
}

// TestTrainSplit -- used during training or exporting
type TestTrainSplit struct {
	TrainImages []string `json:"train_images"`
	TestImages  []string `json:"test_images"`
}

// SSDExport -- for exporting detection with SSD architecture
type SSDExport struct {
	FileName string  `json:"file_name"`
	Width    float32 `json:"width"`
	Height   float32 `json:"height"`
	Class    string  `json:"class"`
	Xmin     float32 `json:"xmin"`
	Ymin     float32 `json:"ymin"`
	Xmax     float32 `json:"xmax"`
	Ymax     float32 `json:"ymax"`
}

// UploadedFile -- contains params captured by multer
type UploadedFile struct {
	FieldName    string  `json:"fieldname"`
	OriginalName string  `json:"originalname"`
	Mimetype     string  `json:"mimetype"`
	Destination  string  `json:"destination"`
	TmpFileName  string  `json:"filename"`
	Path         string  `json:"path"`
	Size         float32 `json:"size"`
	Width        float32 `json:"width"`
	Height       float32 `json:"height"`
}

type AnnInference struct {
	Confidence float64 `json:"confidence"`
	Class      string  `json:"cls"`
}

type ImgInference struct {
	ExperimentImageID bson.ObjectId  `json:"exp_image_id"`
	DatasetImageID    bson.ObjectId  `json:"dataset_image_id"`
	Annotations       []AnnInference `json:"annotations"`
}
