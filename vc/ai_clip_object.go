package vc

import (
	"smcp/models"
	"smcp/utils"
	"strings"
	"time"
)

type AiClipObject struct {
	SourceId         string
	AiQueueModels    []*models.AiClipQueueModel
	FileName         string
	CreatedAt        string
	LastModified     string
	Duration         int
	CreatedAtTime    time.Time
	LastModifiedTime time.Time
}

func (v *AiClipObject) SetupDateTimes() {
	fileName := utils.GetFileNameWithoutExtension(v.FileName)
	v.CreatedAtTime = utils.StringToTime(strings.Split(fileName, ".")[0])
	v.CreatedAt = utils.TimeToString(v.CreatedAtTime, false)
	v.LastModifiedTime = v.CreatedAtTime.Add(time.Duration(v.Duration * int(time.Second)))
	v.LastModified = utils.TimeToString(v.LastModifiedTime, false)
}

func (v *AiClipObject) IsInTimeSpan(check time.Time) bool {
	start := v.CreatedAtTime
	end := v.LastModifiedTime
	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}
	if start.Equal(end) {
		return check.Equal(start)
	}
	return !start.After(check) || !end.Before(check)
}
