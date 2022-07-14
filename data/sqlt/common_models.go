package sqlt

import (
	"smcp/utils"
	"time"
)

type BaseEntity struct {
	GroupId      string `json:"group_id" gorm:"index"`
	SourceId     string `json:"source_id"  gorm:"index:idx_query"` //Index
	CreatedAtStr string `json:"created_at_str"`

	ImageFileName string `json:"image_file_name"`
	VideoFileName string `json:"video_file_name" gorm:"index"` //Index

	AiClipEnabled           bool   `json:"ai_clip_enabled"`
	AiClipFileName          string `json:"ai_clip_file_name" gorm:"index"` //Index
	AiClipCreatedAtStr      string `json:"ai_clip_created_at_str"`
	AiClipLastModifiedAtStr string `json:"ai_clip_last_modified_at_str"`
	AiClipDuration          int    `json:"ai_clip_duration"`

	//extended
	Year   int `json:"year" gorm:"index:idx_query"`  //Index
	Month  int `json:"month" gorm:"index:idx_query"` //Index
	Day    int `json:"day" gorm:"index:idx_query"`   //Index
	Hour   int `json:"hour" gorm:"index:idx_query"`  //Index
	Minute int `json:"minute"`
	Second int `json:"second"`

	CreatedDate time.Time `json:"created_date"` //Index
}

func (b *BaseEntity) SetupDates(createdAt string) {
	b.CreatedAtStr = createdAt
	t := utils.StringToTime(createdAt)
	b.Year = t.Year()
	b.Month = int(t.Month())
	b.Day = t.Day()
	b.Hour = t.Hour()
	b.Minute = t.Minute()
	b.Second = t.Second()
	b.CreatedDate = t
}
