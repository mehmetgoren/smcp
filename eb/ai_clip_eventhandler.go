package eb

import (
	"github.com/go-redis/redis/v8"
	"smcp/reps"
	"smcp/utils"
)

type AiClipEventHandler struct {
	Connection *redis.Client
	AiType     int
}

func (a *AiClipEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()
	rep := reps.AiClipQueueRepository{Connection: a.Connection}
	err := rep.Add(a.AiType, &event.Payload)
	if err != nil {
		return false, err
	}

	return true, nil
}
