package sqlt

import "gorm.io/gorm"

type DetectedPlate struct {
	ImgWidth              int     `json:"img_width"`
	ImgHeight             int     `json:"img_height"`
	TotalProcessingTimeMs float64 `json:"total_processing_time_ms"`

	Plate            string  `json:"plate" gorm:"index:idx_query"` //Index
	Confidence       float64 `json:"confidence"`
	ProcessingTimeMs float64 `json:"processing_time_ms"`

	CandidatesJson string `json:"candidates_json"`
	X0             int    `json:"x0" bson:"x0"`
	Y0             int    `json:"y0" bson:"y0"`
	X1             int    `json:"x1" bson:"x1"`
	Y1             int    `json:"y1" bson:"y1"`
}

type AlprEntity struct {
	gorm.Model
	BaseEntity
	DetectedPlate
}
