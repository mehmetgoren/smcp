package mng

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoFile struct {
	Name            string             `json:"name" bson:"name"` //Index
	CreatedDate     primitive.DateTime `json:"created_date" bson:"created_date"`
	Duration        int                `json:"duration" bson:"duration"`
	Merged          bool               `json:"merged" bson:"merged"`
	ObjectAppearsAt int                `json:"object_appears_at" bson:"object_appears_at"`
}
