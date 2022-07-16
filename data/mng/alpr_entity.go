package mng

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smcp/data"
	"smcp/utils"
	"time"
)

type Candidate struct {
	Plate      string  `json:"plate" bson:"plate"`
	Confidence float64 `json:"confidence" bson:"confidence"`
}

type Coordinates struct {
	X0 int `json:"x0" bson:"x0"`
	Y0 int `json:"y0" bson:"y0"`
	X1 int `json:"x1" bson:"x1"`
	Y1 int `json:"y1" bson:"y1"`
}

type DetectedPlate struct {
	Plate            string  `json:"plate" bson:"plate"` //Index
	Confidence       float64 `json:"confidence" bson:"confidence"`
	ProcessingTimeMs float64 `json:"processing_time_ms" bson:"processing_time_ms"`

	Candidates  []*Candidate `json:"candidates" bson:"candidates"`
	Coordinates *Coordinates `json:"coordinates" bson:"coordinates"`
}

type AlprEntity struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	GroupId   string             `json:"group_id" bson:"group_id"`   //Index
	SourceId  string             `json:"source_id" bson:"source_id"` //Index
	CreatedAt string             `json:"created_at" bson:"created_at"`

	ImgWidth              int            `json:"img_width" bson:"img_width"`
	ImgHeight             int            `json:"img_height" bson:"img_height"`
	TotalProcessingTimeMs float64        `json:"total_processing_time_ms" bson:"total_processing_time_ms"`
	DetectedPlate         *DetectedPlate `json:"detected_plate" bson:"detected_plate"`

	ImageFileName        string     `json:"image_file_name" bson:"image_file_name"`
	VideoFileName        string     `json:"video_file_name" bson:"video_file_name"` //Index
	VideoFileCreatedDate *time.Time `json:"video_file_created_date" bson:"video_file_created_date"`
	VideoFileDuration    int        `json:"video_file_duration" bson:"video_file_duration"`

	AiClip *data.AiClip `json:"ai_clip" bson:"ai_clip"`

	CreatedDate primitive.DateTime `json:"created_date" bson:"created_date"`
}

func (a *AlprEntity) SetupDates(createdAt string) {
	a.CreatedAt = createdAt
	t := utils.StringToTime(createdAt)
	a.CreatedDate = data.TimeToMongoDateTime(t)
}
