package eb

import (
	"github.com/go-redis/redis/v8"
	"log"
	"smcp/data"
	"smcp/data/cmn"
	"smcp/utils"
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

	var vre = &VfiResponseEvent{}
	err := utils.DeserializeJson(event.Payload, vre)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	results := vre.Results
	if results == nil && len(results) == 0 {
		return nil, nil
	}

	rep := v.Factory.CreateRepository()
	for _, r := range results {
		params := &data.SetVideoFileParams{}
		params.Setup(r.SourceId, r.VideoFilename, r.DateStr, r.Duration)
		err = rep.SetVideoFields(params)
	}

	return nil, err
}
