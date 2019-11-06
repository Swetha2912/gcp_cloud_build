package models

// DetectionCoord - mongoDB mapping of annotation coords for detection type
type DetectionCoord struct {
	Class string  `bson:"cls,omitempty" json:"cls,omitempty"`
	X     float32 `bson:"x" json:"x"`
	Y     float32 `bson:"y" json:"y"`
	W     float32 `bson:"w" json:"w"`
	H     float32 `bson:"h" json:"h"`
}

//DetectionJSON -- json format of export segmenattion
type DetectionJSON struct {
	Fileref string `json:"fileref"`
	// FileSize          string                        `json:"size"`
	// FileAttributes        string `json:"file_attributes"`
	Base64ImgData string                     `json:"base64_img_data"`
	Filename      string                     `json:"filename"`
	Regions       map[string]DetectionRegion `json:"regions"`
}

//DetectionRegion -- regions structure
type DetectionRegion struct {
	ShapeAttributes  DetectionShapeAttributes `json:"shape_attributes"`
	RegionAttributes map[string]string        `json:"region_attributes"`
}

//DetectionShapeAttributes -- coordinates
type DetectionShapeAttributes struct {
	Name   string  `json:"name"`
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
}
