package eb

import (
	"github.com/go-redis/redis/v8"
	"smcp/models"
	"smcp/reps"
	"smcp/utils"
)

type FrDiskEventHandler struct {
	Fhr *reps.FrHandlerRepository
}

func (d *FrDiskEventHandler) Handle(event *redis.Message) (interface{}, error) {
	defer utils.HandlePanic()

	var de = models.FaceRecognitionModel{}
	utils.DeserializeJson(event.Payload, &de)
	d.Fhr.Save(&de)

	return nil, nil
}
