package mng

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OdEntityScheme struct {
}

func (o OdEntityScheme) GetCollectionName() string {
	return "od"
}

func (o OdEntityScheme) CreateIndexes(coll *mongo.Collection) ([]string, error) {
	indexes := make([]mongo.IndexModel, 0)

	indexes = append(indexes, mongo.IndexModel{
		Keys: bson.D{
			{Key: "source_id", Value: 1},
			{Key: "created_date", Value: 1},
		},
	})

	indexes = append(indexes, mongo.IndexModel{
		Keys: bson.M{
			"video_file.name": 1, // index in descending order
		},
	})

	indexes = append(indexes, mongo.IndexModel{
		Keys: bson.M{
			"group_id": 1, // index in descending order
		},
	})

	return coll.Indexes().CreateMany(context.TODO(), indexes)
}
