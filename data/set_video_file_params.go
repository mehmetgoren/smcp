package data

import (
	"smcp/utils"
	"time"
)

type SetVideoFileParams struct {
	SourceId      string
	Duration      int
	T1            *time.Time
	T2            *time.Time
	VideoFilename string
}

func (s *SetVideoFileParams) Setup(sourceId string, videoFileName string, dateStr string, duration int) {
	t1 := utils.StringToTime(dateStr)
	t2 := t1.Add(time.Duration(duration) * time.Second)

	s.SourceId = sourceId
	s.VideoFilename = videoFileName
	s.T1 = &t1
	s.T2 = &t2
	s.Duration = duration
}
