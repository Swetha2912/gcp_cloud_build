package models

import (
	"time"
	// "go.mongodb.org/mongo-driver/mongo"
	// "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type YoloTrainingLog struct {
	Step        int       `bson:"step" json:"step"`
	Loss        float64   `bson:"loss" json:"loss"`
	AvgLoss     float64   `bson:"avg_loss" json:"avg_loss"`
	LearnRate   float64   `bson:"learn_rate" json:"learn_rate"`
	TimePerStep float64   `bson:"step_time" json:"step_time"`
	ImagesCount int       `bson:"images_count" json:"images_count"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}

// TrainRequest - MongoDB mapping
type TrainRequest struct {
	ID           bson.ObjectId            `bson:"_id" json:"_id"`
	ExperimentID bson.ObjectId            `bson:"experiment_id,omitempty" json:"experiment_id,omitempty"`
	Classes      []string                 `bson:"classes" json:"classes"`
	Status       string                   `bson:"status" json:"status"`
	Architecture string                   `bson:"architecture" json:"architecture"` // yolo
	Framework    string                   `bson:"framework" json:"framework"`       // tensorflow
	Variant      string                   `bson:"variant" json:"variant"`           // yolov3 or tinyyolov3
	ProcessID    int                      `bson:"process_id" json:"process_id"`
	CreatedAt    time.Time                `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time                `bson:"updated_at" json:"updated_at"`
	Logs         []map[string]interface{} `bson:"logs" json:"logs"`
	Models       []TmpModel               `bson:"models" json:"models"`
	TrainRatio   float64                  `bson:"train_ratio" json:"train_ratio"`
	// yolo specific
	YOLOBatchSize     int    `bson:"yolo_batch_size" json:"yolo_batch_size"`
	YOLOSubdivision   int    `bson:"yolo_subdivision" json:"yolo_subdivision"`
	YOLOModelBoundary string `bson:"yolo_model_boundary" json:"yolo_model_boundary"`
}
