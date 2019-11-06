package models

//SegmentationJSON -- json format of export segmenattion
type SegmentationJSON struct {
	Fileref string `json:"fileref"`
	// FileSize          string                        `json:"size"`
	// FileAttributes        string `json:"file_attributes"`
	Base64ImgData string                        `json:"base64_img_data"`
	Filename      string                        `json:"filename"`
	Regions       map[string]SegmentationRegion `json:"regions"`
}

//SegmentationCSV -- csv format of segmenattion
type SegmentationCSV struct {
	FileName string `json:"#filename"`
	// FileSize string `json:"file_size"`
	// FileAttributes        string `json:"file_attributes"`
	RegionCount           int    `json:"region_count"`
	RegionID              int    `json:"region_id"`
	RegionShapeAttributes string `json:"region_shape_attributes"`
	RegionAttributes      string `json:"region_attributes"`
}

//Region -- region attribute name for CSV format
type Region struct {
	Class string `json:"class"`
}

//SegmentationRegion -- regions structure
type SegmentationRegion struct {
	ShapeAttributes  SegmentationShapeAttributes `json:"shape_attributes"`
	RegionAttributes map[string]string           `json:"region_attributes"`
}

//SegmentationShapeAttributes -- coordinates of the polygon
type SegmentationShapeAttributes struct {
	Name    string    `json:"name"`
	Xpoints []float32 `json:"all_points_x"`
	Ypoints []float32 `json:"all_points_y"`
}

// SegmentationCoord - mongoDB mapping of annotation coords for Segmentation type
type SegmentationCoord struct {
	Class    string      `bson:"cls,omitempty" json:"cls,omitempty"`
	IsActive bool        `bson:"is_active,omitempty" json:"is_active,omitempty"`
	IsSaving bool        `bson:"is_saving,omitempty" json:"is_saving,omitempty"`
	Points   [][]float32 `bson:"points" json:"points"`
}
