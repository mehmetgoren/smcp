package mng

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smcp/data"
	"smcp/utils"
)

type AiDetectedObject struct {
	Label string  `json:"label" bson:"label"`
	Score float32 `json:"score"  bson:"score"`
	X1    float32 `json:"x1" bson:"x1"`
	Y1    float32 `json:"y1" bson:"y1"`
	X2    float32 `json:"x2" bson:"x2"`
	Y2    float32 `json:"y2" bson:"y2"`
}

type AiEntity struct {
	Id     primitive.ObjectID `json:"_id" bson:"_id"`
	Module string             `json:"module" bson:"module"`
	// todo: remove it if it is not necessary
	GroupId        string            `json:"group_id" bson:"group_id"`   //Index
	SourceId       string            `json:"source_id" bson:"source_id"` //Index
	CreatedAt      string            `json:"created_at" bson:"created_at"`
	DetectedObject *AiDetectedObject `json:"detected_object" bson:"detected_object"`
	ImageFileName  string            `json:"image_file_name" bson:"image_file_name"`

	VideoFile *VideoFile `json:"video_file" bson:"video_file"`

	AiClip *data.AiClip `json:"ai_clip" bson:"ai_clip"`

	CreatedDate primitive.DateTime `json:"created_date" bson:"created_date"`
}

func (o *AiEntity) SetupDates(createdAt string) {
	o.CreatedAt = createdAt
	t := utils.StringToTime(createdAt)
	o.CreatedDate = data.TimeToMongoDateTime(t)
}
