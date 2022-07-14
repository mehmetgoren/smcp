package mng

import (
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"smcp/data"
	"smcp/models"
)

type MongoRepository struct {
	Db *DbContext
}

func (o *MongoRepository) OdSave(od *models.ObjectDetectionModel) error {
	m := OdMapper{Config: o.Db.Config}
	entities := m.Map(od)
	return o.Db.Ods.AddRange(entities)
}

func (o *MongoRepository) FrSave(fr *models.FaceRecognitionModel) error {
	m := FrMapper{Config: o.Db.Config}
	entities := m.Map(fr)
	return o.Db.Frs.AddRange(entities)
}

func (o *MongoRepository) AlprSave(alpr *models.AlprResponse) error {
	m := AlprMapper{Config: o.Db.Config}
	entities := m.Map(alpr)
	return o.Db.Alprs.AddRange(entities)
}

func (o *MongoRepository) SetOdVideoClipFields(groupId string, clip *data.AiClip) error {
	if clip == nil {
		return nil
	}
	founds, err := o.Db.Ods.GetByQuery(bson.M{"group_id": groupId})
	if err != nil {
		return err
	}
	if founds != nil && len(founds) > 0 {
		ctx := context.TODO()
		coll := o.Db.Ods.GetCollection()
		for _, entity := range founds {
			entity.AiClip = clip
			filter := bson.M{"_id": entity.Id}
			_, err = coll.ReplaceOne(ctx, filter, entity)
		}
	}

	return err
}
