package data

import (
	"path/filepath"
	"smcp/utils"
	"strings"
	"time"
)

type SetVideoFileMergeParams struct {
	SourceId                string
	OutputFileName          string
	MergedVideoFilenames    []string
	MergedVideoFileDuration int
	VideoFileCreatedDate    time.Time
}

func (s *SetVideoFileMergeParams) Setup(sourceId string, outputFileName string, mergedVideoFileDuration int, mergedVideoFilenames []string) {
	s.SourceId = sourceId
	s.OutputFileName = outputFileName
	s.MergedVideoFilenames = mergedVideoFilenames
	s.MergedVideoFileDuration = mergedVideoFileDuration

	dateStr := strings.Split(filepath.Base(outputFileName), ".")[0]
	s.VideoFileCreatedDate = utils.StringToTime(dateStr)
}

type IVideoFile interface {
	GetEntitiesByName(videoFileName string) ([]interface{}, error)
	GetDuration(entity interface{}) int
	GetEntitiesByNameAndMerged(videoFileName string, merged bool) ([]interface{}, error)
	GetName(entity interface{}) string
	SetObjectAppearsAt(entity interface{}, objectAppearsAt int)
	SetName(entity interface{}, name string)
	SetDuration(entity interface{}, duration int)
	SetMerged(entity interface{}, merged bool)
	SetCreatedDate(entity interface{}, createdDate time.Time)
	Update(entity interface{}) error
}

func GenericVideoFileFunc(vf IVideoFile, params *SetVideoFileMergeParams) error {
	var err error
	prevMergedFileDuration := 0
	prevMergedFiles, err := vf.GetEntitiesByName(params.OutputFileName)
	if prevMergedFiles != nil && len(prevMergedFiles) > 0 {
		prevMergedFileDuration = vf.GetDuration(prevMergedFiles[0])
	}

	entities := make([]interface{}, 0)
	cachedVideoFileDuration := make(map[string]int)
	for _, vfn := range params.MergedVideoFilenames { //Merged Video Files are ordered by created date
		results, err := vf.GetEntitiesByNameAndMerged(vfn, false)
		if err == nil && results != nil && len(results) > 0 {
			for _, entity := range results {
				entities = append(entities, entity)
			}
			if _, ok := cachedVideoFileDuration[vfn]; !ok {
				cachedVideoFileDuration[vfn] = prevMergedFileDuration
				prevMergedFileDuration += vf.GetDuration(results[0])
			}
		}
	}
	for _, entity := range entities {
		vf.SetObjectAppearsAt(entity, cachedVideoFileDuration[vf.GetName(entity)])
		vf.SetName(entity, params.OutputFileName)
		vf.SetDuration(entity, params.MergedVideoFileDuration)
		vf.SetMerged(entity, true)
		vf.SetCreatedDate(entity, params.VideoFileCreatedDate)
		err = vf.Update(entity)
	}

	return err
}
