package sqlt

import (
	"smcp/data"
	"smcp/models"
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

func (s *SqliteRepository) SetOdVideoClipFields(groupId string, clip *data.AiClip) error {
	db := s.Db.Ods.GetGormDb()
	entities := make([]*OdEntity, 0)
	db.Where(map[string]interface{}{"group_id": groupId}).Find(&entities)
	if entities != nil && len(entities) > 0 {
		for _, entity := range entities {
			entity.AiClipEnabled = clip.Enabled
			entity.AiClipFileName = clip.FileName
			entity.AiClipCreatedAtStr = clip.CreatedAt
			entity.AiClipLastModifiedAtStr = clip.LastModifiedAt
			entity.AiClipDuration = clip.Duration

			db.Model(&entity).Updates(&entity)
		}
	}
	return nil
}

func (s *SqliteRepository) SetVideoFileNames(params *data.SetVideoFileNameParams) error {
	db := s.Db.Ods.GetGormDb()
	q := "source_id = ? AND created_date BETWEEN ? AND ?"
	var err error

	updates := map[string]interface{}{"video_file_name": params.VideoFilename, "video_file_created_date": params.T1, "video_file_duration": params.Duration}

	// Od
	ods := make([]*OdEntity, 0)
	db.Where(q, params.SourceId, params.T1, params.T2).Find(&ods)
	for _, od := range ods {
		result := db.Model(od).Where("id=?", od.ID).Updates(updates)
		if result.Error != nil {
			err = result.Error
		}
	}

	//Fr
	frs := make([]*FrEntity, 0)
	db.Where(q, params.SourceId, params.T1, params.T2).Find(&frs)
	for _, fr := range frs {
		result := db.Model(fr).Where("id=?", fr.ID).Updates(updates)
		if result.Error != nil {
			err = result.Error
		}
	}

	//alpr
	alprs := make([]*AlprEntity, 0)
	db.Where(q, params.SourceId, params.T1, params.T2).Find(&alprs)
	for _, alpr := range alprs {
		result := db.Model(alpr).Where("id=?", alpr.ID).Updates(updates)
		if result.Error != nil {
			err = result.Error
		}
	}

	return err
}
