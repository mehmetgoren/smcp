package eb

import (
	"github.com/go-redis/redis/v8"
	"smcp/reps"
	"smcp/utils"
)

type OdAiClipEventHandler struct {
	Connection *redis.Client
}

func (v *OdAiClipEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()
	rep := reps.OdQueueRepository{Connection: v.Connection}
	err := rep.Add(&event.Payload)
	if err != nil {
		return false, err
	}

	return true, nil
}
