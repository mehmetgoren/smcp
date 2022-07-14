package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func TimeToMongoDateTime(time time.Time) primitive.DateTime {
	return primitive.NewDateTimeFromTime(time)
}
