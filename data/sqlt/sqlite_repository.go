package sqlt

import (
	"smcp/data"
	"smcp/models"
	"time"
)

type SqliteRepository struct {
	Db *DbContext
}

func (s *SqliteRepository) OdSave(od *models.ObjectDetectionModel) error {
	m := OdMapper{Config: s.Db.Config}
	entities := m.Map(od)
	return s.Db.Ods.AddRange(entities)
}

func (s *SqliteRepository) FrSave(fr *models.FaceRecognitionModel) error {
	m := FrMapper{Config: s.Db.Config}
	entities := m.Map(fr)
	return s.Db.Frs.AddRange(entities)
}

func (s *SqliteRepository) AlprSave(alpr *models.AlprResponse) error {
	m := AlprMapper{Config: s.Db.Config}
	entities := m.Map(alpr)
	return s.Db.Alprs.AddRange(entities)
}

func (s *SqliteRepository) SetAiClipFields(groupId string, clip *data.AiClip) error {
	//Od
	odDb := s.Db.Ods.GetGormDb()
	odEntities := make([]*OdEntity, 0)
	odDb.Where(map[string]interface{}{"group_id": groupId}).Find(&odEntities)
	if odEntities != nil && len(odEntities) > 0 {
		for _, entity := range odEntities {
			entity.AiClipEnabled = clip.Enabled
			entity.AiClipFileName = clip.FileName
			entity.AiClipCreatedAtStr = clip.CreatedAt
			entity.AiClipLastModifiedAtStr = clip.LastModifiedAt
			entity.AiClipDuration = clip.Duration

			odDb.Model(&entity).Updates(&entity)
		}
	}

	//todo: Be careful, those are not tested.
	//Fr
	frDb := s.Db.Frs.GetGormDb()
	frEntities := make([]*FrEntity, 0)
	frDb.Where(map[string]interface{}{"group_id": groupId}).Find(&frEntities)
	if frEntities != nil && len(frEntities) > 0 {
		for _, entity := range frEntities {
			entity.AiClipEnabled = clip.Enabled
			entity.AiClipFileName = clip.FileName
			entity.AiClipCreatedAtStr = clip.CreatedAt
			entity.AiClipLastModifiedAtStr = clip.LastModifiedAt
			entity.AiClipDuration = clip.Duration

			frDb.Model(&entity).Updates(&entity)
		}
	}

	//todo: Be careful, those are not tested.
	//Alpr
	alprDb := s.Db.Alprs.GetGormDb()
	alprEntities := make([]*AlprEntity, 0)
	alprDb.Where(map[string]interface{}{"group_id": groupId}).Find(&alprEntities)
	if alprEntities != nil && len(alprEntities) > 0 {
		for _, entity := range alprEntities {
			entity.AiClipEnabled = clip.Enabled
			entity.AiClipFileName = clip.FileName
			entity.AiClipCreatedAtStr = clip.CreatedAt
			entity.AiClipLastModifiedAtStr = clip.LastModifiedAt
			entity.AiClipDuration = clip.Duration

			alprDb.Model(&entity).Updates(&entity)
		}
	}

	return nil
}

func createUpdates(params *data.SetVideoFileParams, createdDate *time.Time) map[string]interface{} {
	appearsAtDiff := createdDate.Sub(*params.T1)
	var appearsAt = int(appearsAtDiff.Seconds())
	if appearsAt < 0 {
		appearsAt = 0
	}
	return map[string]interface{}{"video_file_name": params.VideoFilename, "video_file_created_date": params.T1, "video_file_duration": params.Duration,
		"video_file_merged": false, "object_appears_at": appearsAt}
}

func (s *SqliteRepository) SetVideoFields(params *data.SetVideoFileParams) error {
	db := s.Db.Ods.GetGormDb()
	q := "source_id = ? AND created_date BETWEEN ? AND ?"
	var err error

	// Od
	ods := make([]*OdEntity, 0)
	db.Where(q, params.SourceId, params.T1, params.T2).Find(&ods)
	for _, od := range ods {
		result := db.Model(od).Where("id=?", od.ID).Updates(createUpdates(params, &od.CreatedDate))
		if result.Error != nil {
			err = result.Error
		}
	}

	//Fr
	frs := make([]*FrEntity, 0)
	db.Where(q, params.SourceId, params.T1, params.T2).Find(&frs)
	for _, fr := range frs {
		result := db.Model(fr).Where("id=?", fr.ID).Updates(createUpdates(params, &fr.CreatedDate))
		if result.Error != nil {
			err = result.Error
		}
	}

	//Alpr
	alprs := make([]*AlprEntity, 0)
	db.Where(q, params.SourceId, params.T1, params.T2).Find(&alprs)
	for _, alpr := range alprs {
		result := db.Model(alpr).Where("id=?", alpr.ID).Updates(createUpdates(params, &alpr.CreatedDate))
		if result.Error != nil {
			err = result.Error
		}
	}

	return err
}

func (s *SqliteRepository) SetVideoFieldsMerged(params *data.SetVideoFileMergeParams) error {
	err := data.GenericVideoFileFunc(&OdVideoFile{Db: s.Db}, params)
	err = data.GenericVideoFileFunc(&FrVideoFile{Db: s.Db}, params)
	err = data.GenericVideoFileFunc(&AlprVideoFile{Db: s.Db}, params)
	return err
}

//Ods
type OdVideoFile struct {
	Db *DbContext
}

func (o *OdVideoFile) GetEntitiesByName(videoFileName string) ([]interface{}, error) {
	db := o.Db.Ods.GetGormDb()
	entities := make([]*OdEntity, 0)
	err := db.Where("video_file_name = ?", videoFileName).Find(&entities).Error
	return data.TypedToInterfaceArray(entities), err
}

func (o *OdVideoFile) GetDuration(entity interface{}) int {
	od := entity.(*OdEntity)
	return od.VideoFileDuration
}

func (o *OdVideoFile) GetEntitiesByNameAndMerged(videoFileName string, merged bool) ([]interface{}, error) {
	db := o.Db.Ods.GetGormDb()
	entities := make([]*OdEntity, 0)
	err := db.Where("video_file_name = ? AND video_file_merged = ?", videoFileName, merged).Find(&entities).Error
	return data.TypedToInterfaceArray(entities), err
}

func (o *OdVideoFile) GetName(entity interface{}) string {
	od := entity.(*OdEntity)
	return od.VideoFileName
}

func (o *OdVideoFile) GetObjectAppearsAt(entity interface{}) int {
	od := entity.(*OdEntity)
	return od.ObjectAppearsAt
}

func (o *OdVideoFile) GetCreatedDate(entity interface{}) time.Time {
	od := entity.(*OdEntity)
	return od.CreatedDate
}

func (o *OdVideoFile) SetObjectAppearsAt(entity interface{}, objectAppearsAt int) {
	od := entity.(*OdEntity)
	od.ObjectAppearsAt = objectAppearsAt
}

func (o *OdVideoFile) SetName(entity interface{}, name string) {
	od := entity.(*OdEntity)
	od.VideoFileName = name
}

func (o *OdVideoFile) SetDuration(entity interface{}, duration int) {
	od := entity.(*OdEntity)
	od.VideoFileDuration = duration
}

func (o *OdVideoFile) SetMerged(entity interface{}, merged bool) {
	od := entity.(*OdEntity)
	od.VideoFileMerged = merged
}

func (o *OdVideoFile) SetCreatedDate(entity interface{}, createdDate time.Time) {
	od := entity.(*OdEntity)
	od.VideoFileCreatedDate = &createdDate
}

func (o *OdVideoFile) Update(entity interface{}) error {
	od := entity.(*OdEntity)
	return o.Db.Ods.GetGormDb().Model(od).Updates(od).Error
}

//Frs
type FrVideoFile struct {
	Db *DbContext
}

func (f *FrVideoFile) GetEntitiesByName(videoFileName string) ([]interface{}, error) {
	db := f.Db.Frs.GetGormDb()
	entities := make([]*FrEntity, 0)
	err := db.Where("video_file_name = ?", videoFileName).Find(&entities).Error
	return data.TypedToInterfaceArray(entities), err
}

func (f *FrVideoFile) GetDuration(entity interface{}) int {
	fr := entity.(*FrEntity)
	return fr.VideoFileDuration
}

func (f *FrVideoFile) GetEntitiesByNameAndMerged(videoFileName string, merged bool) ([]interface{}, error) {
	db := f.Db.Frs.GetGormDb()
	entities := make([]*FrEntity, 0)
	err := db.Where("video_file_name = ? AND video_file_merged = ?", videoFileName, merged).Find(&entities).Error
	return data.TypedToInterfaceArray(entities), err
}

func (f *FrVideoFile) GetName(entity interface{}) string {
	fr := entity.(*FrEntity)
	return fr.VideoFileName
}

func (f *FrVideoFile) GetObjectAppearsAt(entity interface{}) int {
	fr := entity.(*FrEntity)
	return fr.ObjectAppearsAt
}

func (f *FrVideoFile) GetCreatedDate(entity interface{}) time.Time {
	fr := entity.(*FrEntity)
	return fr.CreatedDate
}

func (f *FrVideoFile) SetObjectAppearsAt(entity interface{}, objectAppearsAt int) {
	fr := entity.(*FrEntity)
	fr.ObjectAppearsAt = objectAppearsAt
}

func (f *FrVideoFile) SetName(entity interface{}, name string) {
	fr := entity.(*FrEntity)
	fr.VideoFileName = name
}

func (f *FrVideoFile) SetDuration(entity interface{}, duration int) {
	fr := entity.(*FrEntity)
	fr.VideoFileDuration = duration
}

func (f *FrVideoFile) SetMerged(entity interface{}, merged bool) {
	fr := entity.(*FrEntity)
	fr.VideoFileMerged = merged
}

func (f *FrVideoFile) SetCreatedDate(entity interface{}, createdDate time.Time) {
	fr := entity.(*FrEntity)
	fr.VideoFileCreatedDate = &createdDate
}

func (f *FrVideoFile) Update(entity interface{}) error {
	fr := entity.(*FrEntity)
	return f.Db.Frs.GetGormDb().Model(fr).Updates(fr).Error
}

//Alpr
type AlprVideoFile struct {
	Db *DbContext
}

func (a *AlprVideoFile) GetEntitiesByName(videoFileName string) ([]interface{}, error) {
	db := a.Db.Alprs.GetGormDb()
	entities := make([]*AlprEntity, 0)
	err := db.Where("video_file_name = ?", videoFileName).Find(&entities).Error
	return data.TypedToInterfaceArray(entities), err
}

func (a *AlprVideoFile) GetDuration(entity interface{}) int {
	alpr := entity.(*AlprEntity)
	return alpr.VideoFileDuration
}

func (a *AlprVideoFile) GetEntitiesByNameAndMerged(videoFileName string, merged bool) ([]interface{}, error) {
	db := a.Db.Alprs.GetGormDb()
	entities := make([]*AlprEntity, 0)
	err := db.Where("video_file_name = ? AND video_file_merged = ?", videoFileName, merged).Find(&entities).Error
	return data.TypedToInterfaceArray(entities), err
}

func (a *AlprVideoFile) GetName(entity interface{}) string {
	alpr := entity.(*AlprEntity)
	return alpr.VideoFileName
}

func (a *AlprVideoFile) GetObjectAppearsAt(entity interface{}) int {
	alpr := entity.(*AlprEntity)
	return alpr.ObjectAppearsAt
}

func (a *AlprVideoFile) GetCreatedDate(entity interface{}) time.Time {
	alpr := entity.(*AlprEntity)
	return alpr.CreatedDate
}

func (a *AlprVideoFile) SetObjectAppearsAt(entity interface{}, objectAppearsAt int) {
	alpr := entity.(*AlprEntity)
	alpr.ObjectAppearsAt = objectAppearsAt
}

func (a *AlprVideoFile) SetName(entity interface{}, name string) {
	alpr := entity.(*AlprEntity)
	alpr.VideoFileName = name
}

func (a *AlprVideoFile) SetDuration(entity interface{}, duration int) {
	alpr := entity.(*AlprEntity)
	alpr.VideoFileDuration = duration
}

func (a *AlprVideoFile) SetMerged(entity interface{}, merged bool) {
	alpr := entity.(*AlprEntity)
	alpr.VideoFileMerged = merged
}

func (a *AlprVideoFile) SetCreatedDate(entity interface{}, createdDate time.Time) {
	alpr := entity.(*AlprEntity)
	alpr.VideoFileCreatedDate = &createdDate
}

func (a *AlprVideoFile) Update(entity interface{}) error {
	alpr := entity.(*AlprEntity)
	return a.Db.Alprs.GetGormDb().Model(alpr).Updates(alpr).Error
}
