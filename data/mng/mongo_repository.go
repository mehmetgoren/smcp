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

	//Ods
	ods, err := o.Db.Ods.GetByQuery(q)
	if err == nil && ods != nil && len(ods) > 0 {
		odColl := o.Db.Ods.GetCollection()
		for _, od := range ods {
			od.VideoFile = createVideoFile(params, od.CreatedDate)
			filter := bson.M{"_id": od.Id}
			_, err = odColl.ReplaceOne(ctx, filter, od)
		}
	}

	//Frs
	frs, err := o.Db.Frs.GetByQuery(q)
	if err == nil && frs != nil && len(frs) > 0 {
		frColl := o.Db.Frs.GetCollection()
		for _, fr := range frs {
			fr.VideoFile = createVideoFile(params, fr.CreatedDate)
			filter := bson.M{"_id": fr.Id}
			_, err = frColl.ReplaceOne(ctx, filter, fr)
		}
	}

	//Alpr
	alprs, err := o.Db.Alprs.GetByQuery(q)
	if err == nil && alprs != nil && len(alprs) > 0 {
		alprColl := o.Db.Alprs.GetCollection()
		for _, alpr := range alprs {
			alpr.VideoFile = createVideoFile(params, alpr.CreatedDate)
			filter := bson.M{"_id": alpr.Id}
			_, err = alprColl.ReplaceOne(ctx, filter, alpr)
		}
	}

	return err
}

func (o *MongoRepository) SetVideoFieldsMerged(params *data.SetVideoFileMergeParams) error {
	err := data.GenericVideoFileFunc(&OdVideoFile{Db: o.Db}, params)
	err = data.GenericVideoFileFunc(&FrVideoFile{Db: o.Db}, params)
	err = data.GenericVideoFileFunc(&AlprVideoFile{Db: o.Db}, params)
	return err
}

// Ods
type OdVideoFile struct {
	Db *DbContext
}

func (o *OdVideoFile) GetEntitiesByName(videoFileName string) ([]interface{}, error) {
	entities, err := o.Db.Ods.GetByQuery(bson.M{"video_file.name": videoFileName})
	return data.TypedToInterfaceArray(entities), err
}

func (o *OdVideoFile) GetDuration(entity interface{}) int {
	od := entity.(*OdEntity)
	return od.VideoFile.Duration
}

func (o *OdVideoFile) GetEntitiesByNameAndMerged(videoFileName string, merged bool) ([]interface{}, error) {
	entities, err := o.Db.Ods.GetByQuery(bson.M{"video_file.name": videoFileName, "video_file.merged": merged})
	return data.TypedToInterfaceArray(entities), err
}

func (o *OdVideoFile) GetName(entity interface{}) string {
	od := entity.(*OdEntity)
	return od.VideoFile.Name
}

func (o *OdVideoFile) SetObjectAppearsAt(entity interface{}, objectAppearsAt int) {
	od := entity.(*OdEntity)
	od.VideoFile.ObjectAppearsAt = objectAppearsAt
}

func (o *OdVideoFile) SetName(entity interface{}, name string) {
	od := entity.(*OdEntity)
	od.VideoFile.Name = name
}

func (o *OdVideoFile) SetDuration(entity interface{}, duration int) {
	od := entity.(*OdEntity)
	od.VideoFile.Duration = duration
}

func (o *OdVideoFile) SetMerged(entity interface{}, merged bool) {
	od := entity.(*OdEntity)
	od.VideoFile.Merged = merged
}

func (o *OdVideoFile) SetCreatedDate(entity interface{}, createdDate time.Time) {
	od := entity.(*OdEntity)
	od.VideoFile.CreatedDate = primitive.NewDateTimeFromTime(createdDate)
}

func (o *OdVideoFile) Update(entity interface{}) error {
	od := entity.(*OdEntity)
	_, err := o.Db.Ods.GetCollection().ReplaceOne(context.TODO(), bson.M{"_id": od.Id}, od)
	return err
}

//Frs
type FrVideoFile struct {
	Db *DbContext
}

func (f *FrVideoFile) GetEntitiesByName(videoFileName string) ([]interface{}, error) {
	entities, err := f.Db.Frs.GetByQuery(bson.M{"video_file.name": videoFileName})
	return data.TypedToInterfaceArray(entities), err
}

func (f *FrVideoFile) GetDuration(entity interface{}) int {
	fr := entity.(*FrEntity)
	return fr.VideoFile.Duration
}

func (f *FrVideoFile) GetEntitiesByNameAndMerged(videoFileName string, merged bool) ([]interface{}, error) {
	entities, err := f.Db.Frs.GetByQuery(bson.M{"video_file.name": videoFileName, "video_file.merged": merged})
	return data.TypedToInterfaceArray(entities), err
}

func (f *FrVideoFile) GetName(entity interface{}) string {
	fr := entity.(*FrEntity)
	return fr.VideoFile.Name
}

func (f *FrVideoFile) SetObjectAppearsAt(entity interface{}, objectAppearsAt int) {
	fr := entity.(*FrEntity)
	fr.VideoFile.ObjectAppearsAt = objectAppearsAt
}

func (f *FrVideoFile) SetName(entity interface{}, name string) {
	fr := entity.(*FrEntity)
	fr.VideoFile.Name = name
}

func (f *FrVideoFile) SetDuration(entity interface{}, duration int) {
	fr := entity.(*FrEntity)
	fr.VideoFile.Duration = duration
}

func (f *FrVideoFile) SetMerged(entity interface{}, merged bool) {
	fr := entity.(*FrEntity)
	fr.VideoFile.Merged = merged
}

func (f *FrVideoFile) SetCreatedDate(entity interface{}, createdDate time.Time) {
	fr := entity.(*FrEntity)
	fr.VideoFile.CreatedDate = primitive.NewDateTimeFromTime(createdDate)
}

func (f *FrVideoFile) Update(entity interface{}) error {
	fr := entity.(*FrEntity)
	_, err := f.Db.Frs.GetCollection().ReplaceOne(context.TODO(), bson.M{"_id": fr.Id}, fr)
	return err
}

//Alprs
type AlprVideoFile struct {
	Db *DbContext
}

func (a *AlprVideoFile) GetEntitiesByName(videoFileName string) ([]interface{}, error) {
	entities, err := a.Db.Alprs.GetByQuery(bson.M{"video_file.name": videoFileName})
	return data.TypedToInterfaceArray(entities), err
}

func (a *AlprVideoFile) GetDuration(entity interface{}) int {
	alpr := entity.(*AlprEntity)
	return alpr.VideoFile.Duration
}

func (a *AlprVideoFile) GetEntitiesByNameAndMerged(videoFileName string, merged bool) ([]interface{}, error) {
	entities, err := a.Db.Alprs.GetByQuery(bson.M{"video_file.name": videoFileName, "video_file.merged": merged})
	return data.TypedToInterfaceArray(entities), err
}

func (a *AlprVideoFile) GetName(entity interface{}) string {
	alpr := entity.(*AlprEntity)
	return alpr.VideoFile.Name
}

func (a *AlprVideoFile) SetObjectAppearsAt(entity interface{}, objectAppearsAt int) {
	alpr := entity.(*AlprEntity)
	alpr.VideoFile.ObjectAppearsAt = objectAppearsAt
}

func (a *AlprVideoFile) SetName(entity interface{}, name string) {
	alpr := entity.(*AlprEntity)
	alpr.VideoFile.Name = name
}

func (a *AlprVideoFile) SetDuration(entity interface{}, duration int) {
	alpr := entity.(*AlprEntity)
	alpr.VideoFile.Duration = duration
}

func (a *AlprVideoFile) SetMerged(entity interface{}, merged bool) {
	alpr := entity.(*AlprEntity)
	alpr.VideoFile.Merged = merged
}

func (a *AlprVideoFile) SetCreatedDate(entity interface{}, createdDate time.Time) {
	alpr := entity.(*AlprEntity)
	alpr.VideoFile.CreatedDate = primitive.NewDateTimeFromTime(createdDate)
}

func (a *AlprVideoFile) Update(entity interface{}) error {
	alpr := entity.(*AlprEntity)
	_, err := a.Db.Alprs.GetCollection().ReplaceOne(context.TODO(), bson.M{"_id": alpr.Id}, alpr)
	return err
}
