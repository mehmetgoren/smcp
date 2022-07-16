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

func (o *MongoRepository) SetVideoFileNames(params *data.SetVideoFileNameParams) error {
	q := bson.M{"source_id": params.SourceId, "created_date": bson.M{"$gte": params.T1, "$lte": params.T2}}
	ctx := context.TODO()

	//Ods
	ods, err := o.Db.Ods.GetByQuery(q)
	if err == nil && ods != nil && len(ods) > 0 {
		odColl := o.Db.Ods.GetCollection()
		for _, od := range ods {
			od.VideoFileName = params.VideoFilename
			od.VideoFileCreatedDate = params.T1
			od.VideoFileDuration = params.Duration
			filter := bson.M{"_id": od.Id}
			_, err = odColl.ReplaceOne(ctx, filter, od)
		}
	}

	//Frs
	frs, err := o.Db.Frs.GetByQuery(q)
	if err == nil && frs != nil && len(frs) > 0 {
		frColl := o.Db.Frs.GetCollection()
		for _, fr := range frs {
			fr.VideoFileName = params.VideoFilename
			fr.VideoFileCreatedDate = params.T1
			fr.VideoFileDuration = params.Duration
			filter := bson.M{"_id": fr.Id}
			_, err = frColl.ReplaceOne(ctx, filter, fr)
		}
	}

	//Alpr
	alprs, err := o.Db.Alprs.GetByQuery(q)
	if err == nil && alprs != nil && len(alprs) > 0 {
		alprColl := o.Db.Alprs.GetCollection()
		for _, alpr := range alprs {
			alpr.VideoFileName = params.VideoFilename
			alpr.VideoFileCreatedDate = params.T1
			alpr.VideoFileDuration = params.Duration
			filter := bson.M{"_id": alpr.Id}
			_, err = alprColl.ReplaceOne(ctx, filter, alpr)
		}
	}

	return err
}
