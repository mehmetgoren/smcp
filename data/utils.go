package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func TimeToMongoDateTime(time time.Time) primitive.DateTime {
	return primitive.NewDateTimeFromTime(time)
}

func TypedToInterfaceArray[T any](entityList []*T) []interface{} {
	var result []interface{}
	for _, entity := range entityList {
		result = append(result, entity)
	}
	return result
}
