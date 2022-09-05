package data

import (
	"path/filepath"
	"smcp/utils"
	"sort"
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
	GetObjectAppearsAt(entity interface{}) int
	GetCreatedDate(entity interface{}) time.Time
	SetObjectAppearsAt(entity interface{}, objectAppearsAt int)
	SetName(entity interface{}, name string)
	SetDuration(entity interface{}, duration int)
	SetMerged(entity interface{}, merged bool)
	SetCreatedDate(entity interface{}, createdDate time.Time)
	Update(entity interface{}) error
}

func GenericVideoFileFunc(vf IVideoFile, params *SetVideoFileMergeParams) error {
	var err error
	if params == nil || params.MergedVideoFilenames == nil || len(params.MergedVideoFilenames) == 0 {
		return err
	}
	index := 0
	prevMergedFileDuration := 0
	prevMergedFiles, err := vf.GetEntitiesByName(params.OutputFileName) //sort yapıp al ,diğer türlü öncekini
	if prevMergedFiles != nil && len(prevMergedFiles) > 0 {
		sort.Slice(prevMergedFiles, func(i, j int) bool {
			return vf.GetCreatedDate(prevMergedFiles[i]).Before(vf.GetCreatedDate(prevMergedFiles[j]))
		})
		prevMergedFileDuration = vf.GetDuration(prevMergedFiles[len(prevMergedFiles)-1])
		index = 1
	}

	entities := make([]interface{}, 0)
	cachedVideoFileDuration := make(map[string]int)
	for ; index < len(params.MergedVideoFilenames); index++ { //Merged Video Files are ordered by created date
		vfn := params.MergedVideoFilenames[index]
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
	sort.Slice(entities, func(i, j int) bool {
		return vf.GetCreatedDate(entities[i]).Before(vf.GetCreatedDate(entities[j]))
	})

	for _, entity := range entities {
		appearsAt := vf.GetDuration(entity) - vf.GetObjectAppearsAt(entity)
		vf.SetObjectAppearsAt(entity, appearsAt+cachedVideoFileDuration[vf.GetName(entity)])
		vf.SetName(entity, params.OutputFileName)
		vf.SetDuration(entity, params.MergedVideoFileDuration)
		vf.SetMerged(entity, true)
		vf.SetCreatedDate(entity, params.VideoFileCreatedDate)
		err = vf.Update(entity)
	}

	return err
}
