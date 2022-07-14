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
