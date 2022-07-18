package eb

import (
	"github.com/go-redis/redis/v8"
	"log"
	"smcp/data"
	"smcp/data/cmn"
	"smcp/utils"
)

type VfmResponseEvent struct {
	SourceId                string   `json:"source_id"`
	OutputFileName          string   `json:"output_file_name"`
	MergedVideoFilenames    []string `json:"merged_video_filenames"`
	MergedVideoFileDuration int      `json:"merged_video_file_duration"`
}

type VfmResponseEventHandler struct {
	Factory *cmn.Factory
}

func (v *VfmResponseEventHandler) Handle(event *redis.Message) (interface{}, error) {
	var evnt VfmResponseEvent
	err := utils.DeserializeJson(event.Payload, &evnt)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if len(evnt.SourceId) == 0 || len(evnt.OutputFileName) == 0 || evnt.MergedVideoFilenames == nil && len(evnt.MergedVideoFilenames) == 0 {
		return nil, nil
	}

	params := &data.SetVideoFileMergeParams{}
	params.Setup(evnt.SourceId, evnt.OutputFileName, evnt.MergedVideoFileDuration, evnt.MergedVideoFilenames)

	rep := v.Factory.CreateRepository()

	return nil, rep.SetVideoFieldsMerged(params)
}
