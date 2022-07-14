package mng

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smcp/data"
	"smcp/utils"
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

	ImageFileName string `json:"image_file_name" bson:"image_file_name"`
	VideoFileName string `json:"video_file_name" bson:"video_file_name"` //Index

	AiClip *data.AiClip `json:"ai_clip" bson:"ai_clip"`

	//extended
	Year   int `json:"year" bson:"year"`   //Index
	Month  int `json:"month" bson:"month"` //Index
	Day    int `json:"day" bson:"day"`     //Index
	Hour   int `json:"hour" bson:"hour"`   //Index
	Minute int `json:"minute" bson:"minute"`
	Second int `json:"second" bson:"second"`

	CreatedDate primitive.DateTime `json:"created_date" bson:"created_date"`
}

func (a *AlprEntity) SetupDates(createdAt string) {
	a.CreatedAt = createdAt
	t := utils.StringToTime(createdAt)
	a.Year = t.Year()
	a.Month = int(t.Month())
	a.Day = t.Day()
	a.Hour = t.Hour()
	a.Minute = t.Minute()
	a.Second = t.Second()
	a.CreatedDate = data.TimeToMongoDateTime(t)
}
