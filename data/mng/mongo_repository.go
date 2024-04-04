package mng

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"smcp/data"
	"smcp/models"
	"time"
)

type MongoRepository struct {
	Db *DbContext
}

func (o *MongoRepository) AiSave(ai *models.AiDetectionModel) error {
	m := AiMapper{Config: o.Db.Config}
	entities := m.Map(ai)
	return o.Db.Ais.AddRange(entities)
}

func (o *MongoRepository) SetAiClipFields(groupId string, clip *data.AiClip) error {
	if clip == nil {
		return nil
	}
	var err error

	aiEntities, err := o.Db.Ais.GetByQuery(bson.M{"group_id": groupId})
	if err != nil {
		return err
	}
	if aiEntities != nil && len(aiEntities) > 0 {
		ctx := context.TODO()
		coll := o.Db.Ais.GetCollection()
		for _, entity := range aiEntities {
			entity.AiClip = clip
			filter := bson.M{"_id": entity.Id}
			_, err = coll.ReplaceOne(ctx, filter, entity)
		}
	}

	return err
}

func createVideoFile(params *data.SetVideoFileParams, createdDate primitive.DateTime) *VideoFile {
	vf := &VideoFile{}
	vf.Name = params.VideoFilename
	vf.CreatedDate = primitive.NewDateTimeFromTime(*params.T1)
	vf.Duration = params.Duration
	vf.Merged = false

	appearsAtDiff := createdDate.Time().Sub(*params.T1)
	var appearsAt = int(appearsAtDiff.Seconds())
	if appearsAt < 0 {
		appearsAt = 0
	}
	vf.ObjectAppearsAt = appearsAt
	return vf
}

func (o *MongoRepository) SetVideoFields(params *data.SetVideoFileParams) error {
	q := bson.M{"source_id": params.SourceId, "created_date": bson.M{"$gte": params.T1, "$lte": params.T2}}
	ctx := context.TODO()

	ais, err := o.Db.Ais.GetByQuery(q)
	if err == nil && ais != nil && len(ais) > 0 {
		aiColl := o.Db.Ais.GetCollection()
		for _, ai := range ais {
			ai.VideoFile = createVideoFile(params, ai.CreatedDate)
			filter := bson.M{"_id": ai.Id}
			_, err = aiColl.ReplaceOne(ctx, filter, ai)
		}
	}

	return err
}

func (o *MongoRepository) SetVideoFieldsMerged(params *data.SetVideoFileMergeParams) error {
	err := data.GenericVideoFileFunc(&AiVideoFile{Db: o.Db}, params)
	return err
}

type AiVideoFile struct {
	Db *DbContext
}

func (o *AiVideoFile) GetEntitiesByName(videoFileName string) ([]interface{}, error) {
	entities, err := o.Db.Ais.GetByQuery(bson.M{"video_file.name": videoFileName})
	return data.TypedToInterfaceArray(entities), err
}

func (o *AiVideoFile) GetDuration(entity interface{}) int {
	ai := entity.(*AiEntity)
	return ai.VideoFile.Duration
}

func (o *AiVideoFile) GetEntitiesByNameAndMerged(videoFileName string, merged bool) ([]interface{}, error) {
	entities, err := o.Db.Ais.GetByQuery(bson.M{"video_file.name": videoFileName, "video_file.merged": merged})
	return data.TypedToInterfaceArray(entities), err
}

func (o *AiVideoFile) GetName(entity interface{}) string {
	ai := entity.(*AiEntity)
	return ai.VideoFile.Name
}

func (o *AiVideoFile) GetObjectAppearsAt(entity interface{}) int {
	ai := entity.(*AiEntity)
	return ai.VideoFile.ObjectAppearsAt
}

func (o *AiVideoFile) GetCreatedDate(entity interface{}) time.Time {
	ai := entity.(*AiEntity)
	return ai.CreatedDate.Time()
}

func (o *AiVideoFile) SetObjectAppearsAt(entity interface{}, objectAppearsAt int) {
	ai := entity.(*AiEntity)
	ai.VideoFile.ObjectAppearsAt = objectAppearsAt
}

func (o *AiVideoFile) SetName(entity interface{}, name string) {
	ai := entity.(*AiEntity)
	ai.VideoFile.Name = name
}

func (o *AiVideoFile) SetDuration(entity interface{}, duration int) {
	ai := entity.(*AiEntity)
	ai.VideoFile.Duration = duration
}

func (o *AiVideoFile) SetMerged(entity interface{}, merged bool) {
	ai := entity.(*AiEntity)
	ai.VideoFile.Merged = merged
}

func (o *AiVideoFile) SetCreatedDate(entity interface{}, createdDate time.Time) {
	ai := entity.(*AiEntity)
	ai.VideoFile.CreatedDate = primitive.NewDateTimeFromTime(createdDate)
}

func (o *AiVideoFile) Update(entity interface{}) error {
	ai := entity.(*AiEntity)
	_, err := o.Db.Ais.GetCollection().ReplaceOne(context.TODO(), bson.M{"_id": ai.Id}, ai)
	return err
}
