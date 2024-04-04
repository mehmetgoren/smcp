package sqlt

import (
	"smcp/data"
	"smcp/models"
	"time"
)

type SqliteRepository struct {
	Db *DbContext
}

func (s *SqliteRepository) AiSave(ai *models.AiDetectionModel) error {
	m := AiMapper{Config: s.Db.Config}
	entities := m.Map(ai)
	return s.Db.Ais.AddRange(entities)
}

func (s *SqliteRepository) SetAiClipFields(groupId string, clip *data.AiClip) error {
	//todo: Be careful, those are not tested.
	aiDb := s.Db.Ais.GetGormDb()
	aiEntities := make([]*AiEntity, 0)
	aiDb.Where(map[string]interface{}{"group_id": groupId}).Find(&aiEntities)
	if aiEntities != nil && len(aiEntities) > 0 {
		for _, entity := range aiEntities {
			entity.AiClipEnabled = clip.Enabled
			entity.AiClipFileName = clip.FileName
			entity.AiClipCreatedAtStr = clip.CreatedAt
			entity.AiClipLastModifiedAtStr = clip.LastModifiedAt
			entity.AiClipDuration = clip.Duration

			aiDb.Model(&entity).Updates(&entity)
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
	db := s.Db.Ais.GetGormDb()
	q := "source_id = ? AND created_date BETWEEN ? AND ?"
	var err error

	ais := make([]*AiEntity, 0)
	db.Where(q, params.SourceId, params.T1, params.T2).Find(&ais)
	for _, ai := range ais {
		result := db.Model(ai).Where("id=?", ai.ID).Updates(createUpdates(params, &ai.CreatedDate))
		if result.Error != nil {
			err = result.Error
		}
	}

	return err
}

func (s *SqliteRepository) SetVideoFieldsMerged(params *data.SetVideoFileMergeParams) error {
	err := data.GenericVideoFileFunc(&AiVideoFile{Db: s.Db}, params)
	return err
}

type AiVideoFile struct {
	Db *DbContext
}

func (o *AiVideoFile) GetEntitiesByName(videoFileName string) ([]interface{}, error) {
	db := o.Db.Ais.GetGormDb()
	entities := make([]*AiEntity, 0)
	err := db.Where("video_file_name = ?", videoFileName).Find(&entities).Error
	return data.TypedToInterfaceArray(entities), err
}

func (o *AiVideoFile) GetDuration(entity interface{}) int {
	ai := entity.(*AiEntity)
	return ai.VideoFileDuration
}

func (o *AiVideoFile) GetEntitiesByNameAndMerged(videoFileName string, merged bool) ([]interface{}, error) {
	db := o.Db.Ais.GetGormDb()
	entities := make([]*AiEntity, 0)
	err := db.Where("video_file_name = ? AND video_file_merged = ?", videoFileName, merged).Find(&entities).Error
	return data.TypedToInterfaceArray(entities), err
}

func (o *AiVideoFile) GetName(entity interface{}) string {
	ai := entity.(*AiEntity)
	return ai.VideoFileName
}

func (o *AiVideoFile) GetObjectAppearsAt(entity interface{}) int {
	ai := entity.(*AiEntity)
	return ai.ObjectAppearsAt
}

func (o *AiVideoFile) GetCreatedDate(entity interface{}) time.Time {
	ai := entity.(*AiEntity)
	return ai.CreatedDate
}

func (o *AiVideoFile) SetObjectAppearsAt(entity interface{}, objectAppearsAt int) {
	ai := entity.(*AiEntity)
	ai.ObjectAppearsAt = objectAppearsAt
}

func (o *AiVideoFile) SetName(entity interface{}, name string) {
	ai := entity.(*AiEntity)
	ai.VideoFileName = name
}

func (o *AiVideoFile) SetDuration(entity interface{}, duration int) {
	ai := entity.(*AiEntity)
	ai.VideoFileDuration = duration
}

func (o *AiVideoFile) SetMerged(entity interface{}, merged bool) {
	ai := entity.(*AiEntity)
	ai.VideoFileMerged = merged
}

func (o *AiVideoFile) SetCreatedDate(entity interface{}, createdDate time.Time) {
	ai := entity.(*AiEntity)
	ai.VideoFileCreatedDate = &createdDate
}

func (o *AiVideoFile) Update(entity interface{}) error {
	ai := entity.(*AiEntity)
	return o.Db.Ais.GetGormDb().Model(ai).Updates(ai).Error
}
