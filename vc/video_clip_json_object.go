package vc

import (
	"smcp/models"
	"smcp/utils"
	"strings"
	"time"
)

type VideoClipJsonObject struct {
	DetectedImage    *models.DetectedImage `json:"detected_image"`
	FileName         string                `json:"file_name"`
	CreatedAt        string                `json:"created_at"`
	LastModified     string                `json:"last_modified"`
	Duration         int                   `json:"duration"`
	CreatedAtTime    time.Time             `json:"-"`
	LastModifiedTime time.Time             `json:"-"`
}

func (v *VideoClipJsonObject) IsInTimeSpan(check time.Time) bool {
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

func (v *VideoClipJsonObject) SetupDateTimes() {
	fileName := utils.GetFileNameWithoutExtension(v.FileName)
	v.CreatedAtTime = utils.StringToTime(strings.Split(fileName, ".")[0], false)
	v.CreatedAt = utils.TimeToString(v.CreatedAtTime, false)
	v.LastModifiedTime = v.CreatedAtTime.Add(time.Duration(v.Duration * int(time.Second)))
	v.LastModified = utils.TimeToString(v.LastModifiedTime, false)
}
