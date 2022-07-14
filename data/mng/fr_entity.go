package mng

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smcp/data"
	"smcp/utils"
)

type DetectedFace struct {
	PredScore   float32 `json:"pred_score" bson:"pred_score"`
	PredClsIdx  int     `json:"pred_cls_idx" bson:"pred_cls_idx"`
	PredClsName string  `json:"pred_cls_name" bson:"pred_cls_name"` //Index
	X1          float32 `json:"x1" bson:"x1"`
	Y1          float32 `json:"y1" bson:"y1"`
	X2          float32 `json:"x2" bson:"x2"`
	Y2          float32 `json:"y2" bson:"y2"`
}

type FrEntity struct {
	Id            primitive.ObjectID `json:"_id" bson:"_id"`
	GroupId       string             `json:"group_id" bson:"group_id"`   //Index
	SourceId      string             `json:"source_id" bson:"source_id"` //Index
	CreatedAt     string             `json:"created_at" bson:"created_at"`
	DetectedFace  *DetectedFace      `json:"detected_face" bson:"detected_face"`
	ImageFileName string             `json:"image_file_name" bson:"image_file_name"`
	VideoFileName string             `json:"video_file_name" bson:"video_file_name"` //Index

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

func (f *FrEntity) SetupDates(createdAt string) {
	f.CreatedAt = createdAt
	t := utils.StringToTime(createdAt)
	f.Year = t.Year()
	f.Month = int(t.Month())
	f.Day = t.Day()
	f.Hour = t.Hour()
	f.Minute = t.Minute()
	f.Second = t.Second()
	f.CreatedDate = data.TimeToMongoDateTime(t)
}
