package eb

import (
	"github.com/go-redis/redis/v8"
	"log"
	"smcp/data"
	"smcp/data/cmn"
	"smcp/utils"
	"time"
)

type ProbeResult struct {
	SourceId      string `json:"source_id"`
	VideoFilename string `json:"video_filename"`
	DateStr       string `json:"date_str"`
	Duration      int    `json:"duration"`
}

type VfiResponseEvent struct {
	Results []*ProbeResult `json:"results"`
}

type VfiResponseEventHandler struct {
	Factory *cmn.Factory
}

func (v *VfiResponseEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()

	var fr = &VfiResponseEvent{}
	err := utils.DeserializeJson(event.Payload, fr)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	results := fr.Results
	if results == nil && len(results) == 0 {
		return nil, nil
	}

	rep := v.Factory.CreateRepository()
	for _, result := range results {
		t1 := utils.StringToTime(result.DateStr)
		t2 := t1.Add(time.Duration(result.Duration) * time.Second)

		params := &data.SetVideoFileNameParams{}
		params.SourceId = result.SourceId
		params.VideoFilename = result.VideoFilename
		params.T1 = &t1
		params.T2 = &t2
		params.Duration = result.Duration

		err = rep.SetVideoFileNames(params)
	}

	return nil, err
}
